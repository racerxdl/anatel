package models

import (
    "github.com/graphql-go/graphql"
    "github.com/jinzhu/gorm"
    "time"
)

type CallSign struct {
	gorm.Model
	Callsign    string `gorm:"primary_key"`
	FirstSaw    time.Time
	LastUpdated time.Time
	Class       string
	Region      string
	Stations    []StationData         `gorm:"foreignkey:Callsign;association_foreignkey:Callsign"`
	Repeaters   []RepeaterStationData `gorm:"foreignkey:Callsign;association_foreignkey:Callsign"`
}

func CallSignArrayToString(data []CallSign) []string {
	arr := make([]string, 0)
	for i := 0; i < len(data); i++ {
		arr = append(arr, data[i].Callsign)
	}
	return arr
}

var GQLCallSign = graphql.NewObject(graphql.ObjectConfig{
	Name: "Callsign",
	Fields: graphql.Fields{
		"Callsign": &graphql.Field{Type: graphql.String},
		"FirstSaw": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(CallSign).FirstSaw.String(), nil
			},
		},
		"LastUpdated": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(CallSign).LastUpdated.String(), nil
			},
		},
		"Class":     &graphql.Field{Type: graphql.String},
		"Region":    &graphql.Field{Type: graphql.String},
		"Stations":  &graphql.Field{Type: graphql.NewList(GQLStation)},
		"Repeaters": &graphql.Field{Type: graphql.NewList(GQLRepeaterStation)},
	},
})
