package schemas

import (
	"github.com/graphql-go/graphql"
)

// Define Metrics Type
var metricsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Metrics",
	Fields: graphql.Fields{
		"timeStarted": &graphql.Field{
			Type: graphql.String,
		},
		"timeEnded": &graphql.Field{
			Type: graphql.String,
		},
		"timeStamps": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"cpuUnits": &graphql.Field{
			Type: graphql.String,
		},
		"memoryUnits": &graphql.Field{
			Type: graphql.String,
		},
		"pods": &graphql.Field{
			Type: graphql.NewList(podType),
		},
		"nodes": &graphql.Field{
			Type: graphql.NewList(nodeType),
		},
		"grouped": &graphql.Field{
			Type: graphql.NewList(groupedType),
		},
	},
})

// Define Pod Type with lists for metric data
var podType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Pod",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"cpuRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"memoryRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"cpuHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"memoryHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})

// Define Node Type with lists for metric data
var nodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Node",
	Fields: graphql.Fields{
		"nodeName": &graphql.Field{
			Type: graphql.String,
		},
		"cpuRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"memoryRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"cpuPercentage": &graphql.Field{
			Type: graphql.NewList(graphql.Float),
		},
		"memoryPercentage": &graphql.Field{
			Type: graphql.NewList(graphql.Float),
		},
		"cpuHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"memoryHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"cpuPercentageHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"memoryPercentageHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})

// Define Grouped Type with lists for metric data (similar to Pod)
var groupedType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Grouped",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"cpuRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"memoryRaw": &graphql.Field{
			Type: graphql.NewList(graphql.Int),
		},
		"cpuHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"memoryHumanReadable": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
	},
})
