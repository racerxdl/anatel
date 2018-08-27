package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"golang.org/x/crypto/sha3"
	"strconv"
	"strings"
	"encoding/hex"
	"github.com/quan-to/graphql"
)

type StationData struct {
	gorm.Model
	UID string
	Owner string
	Certificate string
	InclusionDate time.Time
	ExpirationDate time.Time
	Callsign string
	City string
	Service string
	StationType string
	Region string
	FirstSaw time.Time
}

func (cs *StationData) GenerateUID() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strings.ToLower(cs.Callsign)))
	hasher.Write([]byte(strconv.FormatInt(cs.InclusionDate.UnixNano(), 10)))
	hasher.Write([]byte(strconv.FormatInt(cs.ExpirationDate.UnixNano(), 10)))
	hasher.Write([]byte(strings.ToLower(cs.StationType)))
	hasher.Write([]byte(strings.ToLower(cs.Owner)))

	result = hasher.Sum(result)

	cs.UID = strings.ToUpper(hex.EncodeToString(result))
}

var GQLStation = graphql.NewObject(graphql.ObjectConfig{
	Name: "StationData",
	Fields: graphql.Fields{
		"UID": &graphql.Field{ Type: graphql.String },
		"Owner": &graphql.Field{ Type: graphql.String },
		"InclusionDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(StationData).InclusionDate.String(), nil
			},
		},
		"Certificate": &graphql.Field{ Type: graphql.String },
		"ExpirationDate": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(StationData).ExpirationDate.String(), nil
			},
		},
		"Callsign": &graphql.Field{ Type: graphql.String },
		"Service": &graphql.Field{ Type: graphql.String },
		"City": &graphql.Field{ Type: graphql.String },
		"StationType": &graphql.Field{ Type: graphql.String },
		"Region": &graphql.Field{ Type: graphql.String },
		"FirstSaw": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(StationData).FirstSaw.String(), nil
			},
		},
	},
})
