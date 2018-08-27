package main

import (
	"github.com/quan-to/handler"
	"github.com/quan-to/graphql"
	"log"
	"net/http"
	"github.com/racerxdl/anatel/gql"
	"github.com/racerxdl/anatel/models"
	"github.com/quan-to/graphql/gqlerrors"
	"runtime/debug"
)

func graphqlServer() {
	var rootQuery = graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"Callsign": &graphql.Field{
			Type: gql.MakeGraphQLConnection(models.GQLCallSign),
			Description: "Callsign Query",
			Args: gql.MakeConnectionArgs(graphql.FieldConfigArgument{
				"Callsign": {
					Type: graphql.String,
					Description: "Callsign to search (can be partial, empty to return all)",
				},
				"Region": {
					Type: graphql.String,
					Description: "Region to search (a.k.a. UF)",
				},
			}),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return SearchCallsigns(p.Args, database), nil
			},
		},
		"Repeater": &graphql.Field{
			Type: gql.MakeGraphQLConnection(models.GQLRepeaterStation),
			Description: "Repeaters Query",
			Args: gql.MakeConnectionArgs(graphql.FieldConfigArgument{
				"Callsign": {
					Type: graphql.String,
					Description: "Callsign to search (can be partial, empty to return all)",
				},
				"Region": {
					Type: graphql.String,
					Description: "Region to search (a.k.a. UF)",
				},
			}),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return SearchRepeater(p.Args, database), nil
			},
		},
	}}

	var schemaConfig = graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: false,
		Playground: true,
		CustomErrorFormatter: func(err error) gqlerrors.FormattedError {
			log.Println(err)
			log.Println(string(debug.Stack()))
			return gqlerrors.FormatError(err)
		},
	})

	http.Handle("/graphql", h)
	log.Println("Listening on :5000")
	http.ListenAndServe(":5000", nil)
	log.Println("Exiting")
}
