package schemas

import "github.com/graphql-go/graphql"

// TraceQueryInputType Define an input object type for the "traces" field arguments
var TraceQueryInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "TraceQueryInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"ids": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(graphql.String), // List of trace IDs
		},
		"statusCode": &graphql.InputObjectFieldConfig{
			Type: graphql.Int, // Status code filter
		},
		"beginTime": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // Begin time filter
		},
		"endTime": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // End time filter
		},
		"totalTimeHigherThan": &graphql.InputObjectFieldConfig{
			Type: graphql.Int, // Total time filter
		},
		"podName": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // podName filter
		},
		"operation": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // podName filter
		},
	},
})

var MetricsQueryInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MetricsQueryInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"beginTime": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // Begin time filter
		},
		"endTime": &graphql.InputObjectFieldConfig{
			Type: graphql.String, // End time filter
		},
		"order": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})
