package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"os"
	"github.com/racerxdl/anatel/models"
	"log"
	"strconv"
	"encoding/base64"
	"github.com/racerxdl/anatel/gql"
)

func Initialize() *gorm.DB {
	db, err := gorm.Open("postgres",
		fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		))

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.CallSign{})
	db.AutoMigrate(&models.CertificateData{})
	db.AutoMigrate(&models.StationData{})
	db.AutoMigrate(&models.TestData{})
	db.AutoMigrate(&models.RepeaterStationData{})

	return db
}


func WriteCallSigns(data []models.CallSign, db *gorm.DB) []int {
	addedCallsigns := make([]int, 0)
	for i := 0; i < len(data); i++ {
		var mcl = data[i]
		var count int
		db.Model(&models.CallSign{}).Where("callsign = ?", mcl.Callsign).Count(&count)

		if count == 0 {
			//log.Printf("Adding %s to the database.\n", mcl.Callsign)
			db.NewRecord(mcl)
			db.Create(&mcl)
			addedCallsigns = append(addedCallsigns, i)
		}

	}

	log.Printf("Added %d new callsigns to database.\n", len(addedCallsigns))
	return addedCallsigns
}

func WriteStationData(data []models.StationData, db *gorm.DB) []int {
	addedStations := make([]int, 0)
	for i := 0; i < len(data); i++ {
		var mcl = data[i]
		var count int
		db.Model(&models.StationData{}).Where("UID = ?", mcl.UID).Count(&count)

		if count == 0 {
			log.Printf("Adding %s (%d) to the database.\n", mcl.Callsign, i)
			db.NewRecord(mcl)
			db.Create(&mcl)
			addedStations = append(addedStations, i)
		}

	}

	return addedStations
}


func WriteRepeaterData(data []models.RepeaterStationData, db *gorm.DB) []int {
	addedStations := make([]int, 0)
	for i := 0; i < len(data); i++ {
		var mcl = data[i]
		var count int
		db.Model(&models.RepeaterStationData{}).Where("UID = ?", mcl.UID).Count(&count)

		if count == 0 {
			log.Printf("Adding Repeater %s (%d) to the database.\n", mcl.Callsign, i)
			db.NewRecord(mcl)
			db.Create(&mcl)
			addedStations = append(addedStations, i)
		}

	}

	return addedStations
}

func WriteTests(data []models.TestData, db *gorm.DB) []int {
	addedTests := make([]int, 0)
	for i := 0; i < len(data); i++ {
		var mcl = data[i]
		var count int
		db.Model(&models.TestData{}).Where("uid = ?", mcl.UID).Count(&count)

		if count == 0 {
			log.Printf("Adding %s to the database.\n", mcl.UID)
			db.NewRecord(mcl)
			db.Create(&mcl)
			addedTests = append(addedTests, i)
		}
		// TODO: Update Status
	}

	return addedTests
}

func FromGormToCursor(id uint) string {
	var sid = strconv.FormatInt(int64(id), 10)
	return base64.StdEncoding.EncodeToString([]byte(sid))
}

func FromCursorToGorm(cursor string) int64 {
	var sid, err = base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		panic(err)
	}

	id, err := strconv.ParseInt(string(sid), 10, 64)
	if err != nil {
		panic(err)
	}

	return id
}

func SearchCallsigns(args map[string]interface{}, db *gorm.DB) (output *gql.ConnectionData) {
	var uf, callsign string
	first := 10
	last := -1
	after := int64(-1)
	before := int64(-1)

	if val, ok := args["Region"] ; ok {
		uf = val.(string)
	}

	if val, ok := args["Callsign"] ; ok {
		callsign = val.(string)
	}

	if val, ok := args["After"] ; ok {
		after = FromCursorToGorm(val.(string))
	}

	if val, ok := args["Before"] ; ok {
		before = FromCursorToGorm(val.(string))
	}

	if val, ok := args["First"] ; ok {
		first = val.(int)
	}

	if val, ok := args["Last"] ; ok {
		last = val.(int)
	}

	if last != -1 {
		panic("Last not implemented")
	}

	if first > 20 {
		first = 20
	}

	s := db.Model(&models.CallSign{})

	if uf != "" {
		s = s.Where("region = ?", uf)
	}

	if callsign != "" {
		s = s.Where("callsign LIKE ?", fmt.Sprintf("%%%s%%", callsign))
	}

	var totalCount int64
	s.Count(&totalCount)

	if after != -1 {
		s = s.Where("id >= ?", after)
	}

	if before != -1 {
		s = s.Where("id <= ?", before)
	}

	s = s.Limit(first)

	s = s.Preload("Stations").Preload("Repeaters")

	var nodes []models.CallSign
	var startCursor, endCursor string

	s.Find(&nodes)

	if len(nodes) > 0 {
		v := nodes[0]
		startCursor = FromGormToCursor(v.ID)
		v = nodes[len(nodes)-1]
		endCursor = FromGormToCursor(v.ID)
	}

	gNodes := make([]interface{}, len(nodes))
	for i, v := range nodes {
		gNodes[i] = v
	}

	var edges = gql.MakeEdges(gNodes, func(m interface{}) string {
		var z = m.(models.CallSign)
		return FromGormToCursor(z.ID)
	})

	output = &gql.ConnectionData{
		TotalCount: totalCount,
		PageInfo: gql.PageInfo{
			StartCursor: startCursor,
			EndCursor: endCursor,
		},
		Edges: edges,
	}

	return output
}

func SearchRepeater(args map[string]interface{}, db *gorm.DB) (output *gql.ConnectionData) {
	var uf, callsign string
	first := 10
	last := -1
	after := int64(-1)
	before := int64(-1)

	if val, ok := args["Region"] ; ok {
		uf = val.(string)
	}

	if val, ok := args["Callsign"] ; ok {
		callsign = val.(string)
	}

	if val, ok := args["After"] ; ok {
		after = FromCursorToGorm(val.(string))
	}

	if val, ok := args["Before"] ; ok {
		before = FromCursorToGorm(val.(string))
	}

	if val, ok := args["First"] ; ok {
		first = val.(int)
	}

	if val, ok := args["Last"] ; ok {
		last = val.(int)
	}

	if last != -1 {
		panic("Last not implemented")
	}

	if first > 20 {
		first = 20
	}

	s := db.Model(&models.RepeaterStationData{})

	if uf != "" {
		s = s.Where("region = ?", uf)
	}

	if callsign != "" {
		s = s.Where("callsign LIKE ?", fmt.Sprintf("%%%s%%", callsign))
	}

	var totalCount int64
	s.Count(&totalCount)

	if after != -1 {
		s = s.Where("id >= ?", after)
	}

	if before != -1 {
		s = s.Where("id <= ?", before)
	}

	s = s.Limit(first)

	var nodes []models.RepeaterStationData
	var startCursor, endCursor string

	s.Find(&nodes)

	if len(nodes) > 0 {
		v := nodes[0]
		startCursor = FromGormToCursor(v.ID)
		v = nodes[len(nodes)-1]
		endCursor = FromGormToCursor(v.ID)
	}

	gNodes := make([]interface{}, len(nodes))
	for i, v := range nodes {
		gNodes[i] = v
	}

	var edges = gql.MakeEdges(gNodes, func(m interface{}) string {
		var z = m.(models.RepeaterStationData)
		return FromGormToCursor(z.ID)
	})

	output = &gql.ConnectionData{
		TotalCount: totalCount,
		PageInfo: gql.PageInfo{
			StartCursor: startCursor,
			EndCursor: endCursor,
		},
		Edges: edges,
	}

	return output
}

