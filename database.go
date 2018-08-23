package main

import "github.com/jinzhu/gorm"
import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"os"
	"github.com/racerxdl/anatel/models"
	"log"
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