package schemas

import (
	"github.com/graphql-go/graphql"
)

var ParentTrace = graphql.NewObject(graphql.ObjectConfig{
	Name: "ParentTrace",
	Fields: graphql.Fields{
		"traceID": &graphql.Field{
			Type: graphql.String,
		},
		"isActive": &graphql.Field{
			Type: graphql.Boolean,
		},
		"items": &graphql.Field{
			Type: graphql.NewList(itemType),
		},
		"timeEnded": &graphql.Field{
			Type: graphql.String,
		},
		"timeStarted": &graphql.Field{
			Type: graphql.String,
		},
		"totalTime": &graphql.Field{
			Type: graphql.Int,
		},
		"responseCode": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var itemType = graphql.NewUnion(graphql.UnionConfig{
	Name: "Item",
	ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
		// Implement logic to determine the type of the object based on the provided data.
		// You can inspect the data and return the appropriate GraphQL type.
		// For example, you can check a field in the data or use other criteria to determine the type.
		// Return the corresponding GraphQL object type.

		// Example:
		if _, isDatabaseSpan := p.Value.(*DatabaseSpan); isDatabaseSpan {
			return databaseSpanType
		} else if _, isSpan := p.Value.(*Span); isSpan {
			return spanType
		} else if _, isTrace := p.Value.(*Trace); isTrace {
			return traceType
		}

		return nil // Return nil if the type cannot be determined.
	},
	Types: []*graphql.Object{
		databaseSpanType,
		spanType,
		traceType,
	},
})

// Define a GraphQL object type for "databaseSpan"
var databaseSpanType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DatabaseSpan",
	Fields: graphql.Fields{
		"parentSpanID": &graphql.Field{
			Type: graphql.String,
		},
		"spanID": &graphql.Field{
			Type: graphql.String,
		},
		"itemType": &graphql.Field{
			Type: graphql.String,
		},
		"query": &graphql.Field{
			Type: graphql.String,
		},
		"namespace": &graphql.Field{
			Type: graphql.String,
		},
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
		"podName": &graphql.Field{
			Type: graphql.String,
		},
		"took": &graphql.Field{
			Type: graphql.String,
		},
		"hits": &graphql.Field{
			Type: graphql.Int,
		},
		"timeFinished": &graphql.Field{
			Type: graphql.String,
		},
		"timeStarted": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Define a GraphQL object type for "span"
var spanType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Span",
	Fields: graphql.Fields{
		"parentSpanID": &graphql.Field{
			Type: graphql.String,
		},
		"requestBody": &graphql.Field{
			Type: graphql.String,
		},
		"spanID": &graphql.Field{
			Type: graphql.String,
		},
		"itemType": &graphql.Field{
			Type: graphql.String,
		},
		"namespace": &graphql.Field{
			Type: graphql.String,
		},
		"action": &graphql.Field{
			Type: graphql.String,
		},
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
		"podName": &graphql.Field{
			Type: graphql.String,
		},
		"timeFinished": &graphql.Field{
			Type: graphql.String,
		},
		"timeStarted": &graphql.Field{
			Type: graphql.String,
		},
		"responseCode": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

// Define a GraphQL object type for "trace"
var traceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Trace",
	Fields: graphql.Fields{
		"parentSpanID": &graphql.Field{
			Type: graphql.String,
		},
		"method": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"host": &graphql.Field{
			Type: graphql.String,
		},
		"remoteAddress": &graphql.Field{
			Type: graphql.String,
		},
		"timestamp": &graphql.Field{
			Type: graphql.String,
		},
		"podName": &graphql.Field{
			Type: graphql.String,
		},
		"namespace": &graphql.Field{
			Type: graphql.String,
		},
		"itemType": &graphql.Field{
			Type: graphql.String,
		},
		"operation": &graphql.Field{
			Type: graphql.String,
		},
		"rootQuery": &graphql.Field{
			Type: graphql.String,
		},
		"metrics": &graphql.Field{
			Type: tracingMetrics,
		},
	},
})

var tracingMetrics = graphql.NewObject(graphql.ObjectConfig{
	Name: "tracingMetrics",
	Fields: graphql.Fields{
		"cpuRaw": &graphql.Field{
			Type: graphql.Int,
		},
		"memoryRaw": &graphql.Field{
			Type: graphql.Int,
		},
		"cpuHumanReadable": &graphql.Field{
			Type: graphql.String,
		},
		"memoryHumanReadable": &graphql.Field{
			Type: graphql.String,
		},
		"memoryUnits": &graphql.Field{
			Type: graphql.String,
		},
		"cpuUnits": &graphql.Field{
			Type: graphql.String,
		},
	},
})
