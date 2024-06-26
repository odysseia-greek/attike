package schemas

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/odysseia-greek/attike/euripides/handlers"
	"log"
	"os"
	"sync"
)

var (
	handler              *handlers.EuripidesHandler
	EuripidesHandlerOnce sync.Once
)

func InitEuripidesHandler() {
	EuripidesHandlerOnce.Do(func() {
		env := os.Getenv("ENV")
		euripidesHandler, err := handlers.CreateNewConfig(env)
		if err != nil {
			log.Print(err)
		}
		handler = euripidesHandler
	})
}

var EuripidesSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		// traces
		"traces": &graphql.Field{
			Type:        graphql.NewList(ParentTrace),
			Description: "Get traces",
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{
					Type: TraceQueryInputType, // Use the custom input object type
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Retrieve the input argument
				input, _ := p.Args["input"].(map[string]interface{})

				// Use "ids" and "queryAll" to customize your tracing query
				// Replace this with your actual resolver logic
				result, err := handler.Tracing(input)
				if err != nil {
					return nil, err
				}

				res := parseTracesToGraphql(result)
				return res, nil
			},
		},
		"metrics": &graphql.Field{
			Type: metricsType,
			Args: graphql.FieldConfigArgument{
				"timeSpan": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"order": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				timeSpan, isOK := p.Args["timeSpan"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument timeSpan")
				}
				order, isOK := p.Args["order"].(string)
				if !isOK {
					return nil, fmt.Errorf("expected argument order")
				}
				result, err := handler.Metrics(timeSpan, order)
				if err != nil {
					return nil, err
				}

				res := parseMetricsToGraphql(result)
				return res, nil
			},
		},
	},
})
