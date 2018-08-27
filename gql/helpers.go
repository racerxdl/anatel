package gql

import (
	"github.com/quan-to/graphql"
	"fmt"
	"log"
)

type PageInfo struct {
	StartCursor string
	EndCursor string
	HasNextPage bool
	HasPreviousPage bool
}

type EdgeData struct {
	Node interface{}
	Cursor string
}

type ConnectionData struct {
	TotalCount int64
	PageInfo PageInfo
	Edges []EdgeData
}

func MakeConnectionArgs(arguments graphql.FieldConfigArgument) graphql.FieldConfigArgument {
	var out = graphql.FieldConfigArgument{
		"First": {
			Type: graphql.Int,
			Description: "Show first N results starting from cursor",
		},
		"Last": {
			Type: graphql.Int,
			Description: "Show last N results starting from cursor",
		},
		"After": {
			Type: graphql.String,
			Description: "Show results after specified cursor",
		},
		"Before": {
			Type: graphql.String,
			Description: "Show results before specified cursor",
		},
	}

	for k, v := range arguments {
		if _, ok := out[k]; ok {
			log.Printf("WARN: There is already a field named %s in argument. Skiping adding it.", k)
		} else {
			out[k] = v
		}
	}

	return out
}

var pageInfoModel = graphql.NewObject(graphql.ObjectConfig{
	Name: "EdgePageInfo",
	Fields: graphql.Fields{
		"StartCursor": {
			Type: graphql.String,
			Description: "Cursor of the first page item",
		},
		"EndCursor": {
			Type: graphql.String,
			Description: "Cursor of the last page item",
		},
		"HasNextPage": {
			Type: graphql.Boolean,
			Description: "If there is a next page",
		},
		"HasPreviousPage": {
			Type: graphql.Boolean,
			Description: "If there is a previous page",
		},
	},
})

func makeGraphQLEdge(nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "GraphQLConnectionEdges",
		Description: "Edges for GraphQL Pagination",
		Fields: graphql.Fields{
			"Node": {
				Type: nodeType,
				Description: fmt.Sprintf("Object %s", nodeType.Name()),
			},
			"Cursor": {
				Type: graphql.String,
				Description: "Cursor of the row",
			},
		},
	})
}

func MakeGraphQLConnection(nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: fmt.Sprintf("%sConnection", nodeType.Name()),
		Description: fmt.Sprintf("Connections for model %s", nodeType.Name()),
		Fields: graphql.Fields{
			"TotalCount": {
				Type: graphql.Int,
				Description: "Total Count of Objects in Query",
			},
			"PageInfo": {
				Type: pageInfoModel,
				Description: "Pagination Information",
			},
			"Edges": {
				Type: graphql.NewList(makeGraphQLEdge(nodeType)),
				Description: "Rows of the query",
			},
		},
	})
}


func MakeEdges(nodes []interface{}, cursorFunc func(m interface{}) string) []EdgeData {
	var ret = make([]EdgeData, len(nodes))

	for i := 0; i < len(nodes); i++ {
		ret[i] = EdgeData{
			Node: nodes[i],
			Cursor: cursorFunc(nodes[i]),
		}
	}

	return ret
}