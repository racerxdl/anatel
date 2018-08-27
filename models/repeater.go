package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"github.com/quan-to/graphql"
)

type RepeaterStationData struct {
	gorm.Model
	UID string
	RXFrequency uint64
	TXFrequency uint64
	Callsign string
	StationNumber string
	City string
	Region string
	StationType string
	FirstSaw time.Time
}

func (cs *RepeaterStationData) GenerateUID() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strings.ToLower(cs.Callsign)))
	hasher.Write([]byte(strings.ToLower(cs.StationNumber)))
	hasher.Write([]byte(strings.ToLower(cs.StationType)))

	result = hasher.Sum(result)

	cs.UID = strings.ToUpper(hex.EncodeToString(result))
}

var GQLRepeaterStation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RepeaterStationData",
	Fields: graphql.Fields{
		"UID": &graphql.Field{ Type: graphql.String },
		"RXFrequency": &graphql.Field{ Type: graphql.Float },
		"TXFrequency": &graphql.Field{ Type: graphql.Float },
		"Callsign": &graphql.Field{ Type: graphql.String },
		"StationNumber": &graphql.Field{ Type: graphql.String },
		"Region": &graphql.Field{ Type: graphql.String },
		"City": &graphql.Field{ Type: graphql.String },
		"StationType": &graphql.Field{ Type: graphql.String },
		"FirstSaw": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(RepeaterStationData).FirstSaw.String(), nil
			},
		},
	},
})
