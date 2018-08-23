package main

import (
	"log"
	"github.com/racerxdl/anatel/models"
	"github.com/tebeka/selenium"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

func UpdateCallSigns(username, password, uf, class string, db *gorm.DB, driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered", r)
		}
	}()
	driver.DeleteAllCookies()
	log.Println("Fetching callsigns")
	callsignsRaw := consultaIndicativos(username, password, uf, class, driver)

	log.Println("Santizing Callsigns")
	callSigns := models.MapCallsignRawData(callsignsRaw)

	log.Println("Writting Callsigns")
	WriteCallSigns(callSigns, db)
}

func updateStationsSegment(length int, state string, firstTime time.Time, db *gorm.DB, driver selenium.WebDriver) {
	var cls []models.CallSign
	// No need for start since we wait for saving callsigns
	db.Model(&models.CallSign{}).Where("last_updated <= ? and region = ? ", firstTime, state).Limit(length).Find(&cls)

	if len(cls) <= 0 {
		return
	}

	log.Println("Fetching Extended Data")

	extDataRaw := consultaIndicativoArray(models.CallSignArrayToString(cls), driver)
	extData := models.MapStationRawData(extDataRaw)

	ns := make(chan []int)

	go func() {
		log.Println("Writting Station Data")
		newStations := WriteStationData(extData, db)

		log.Println("Updating Callsigns")
		for i := 0; i < len(cls); i++ {
			db.Model(&cls[i]).Update("last_updated", time.Now())
		}
		ns <- newStations
	}()

	log.Println("Waiting for Station / Callsign data to be written")

	newStations := <- ns

	log.Println("Triggering Notifications")

	triggerStationNotifications(newStations, extData, db)
}

func UpdateStationsFlow(state string, db *gorm.DB, driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered", r)
		}
	}()
	driver.DeleteAllCookies()
	log.Printf("Fetching callsigns for region %s\n", state)

	var firstTime = time.Now().Add(-callSignUpdateTimeout)

	// Grab Count
	var numCallSign int
	var totalCallsigns int
	db.Model(&models.CallSign{}).Where("last_updated <= ? and region = ?", firstTime, state).Count(&numCallSign)
	totalCallsigns = numCallSign
	var numSegments = numCallSign / segmentLength

	if numSegments == 0 {
		numSegments = 1
	}

	var pos = 0

	for i := 0; i < numSegments; i++ {
		var chunkSize = segmentLength
		if chunkSize > numCallSign {
			chunkSize = numCallSign
		}

		log.Printf("Fetching callsigns from %d to %d from %d records\n", pos, pos + chunkSize, totalCallsigns)
		updateStationsSegment(chunkSize, state, firstTime, db, driver)
		numCallSign -= chunkSize
	}
}

const testCheckMonths = 6

func GetNextTests(username, password, uf string, db *gorm.DB, driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error:", r)
		}
	}()

	driver.DeleteAllCookies()
	startTime := time.Now() // .Add(-time.Hour * 24 * 30)

	//startTime = startTime.Add(-time.Hour * 24 * time.Duration(startTime.Day() - 1)) // Reset for the first day of month

	testsRaw := make([]map[string]string, 0)

	skipLogin := false

	for i := 0; i < testCheckMonths; i++ {
		v := consultaAgenda(username, password, uf, startTime, skipLogin, driver)
		testsRaw = append(testsRaw, v...)
		startTime = startTime.Add(time.Hour * 24 * 30 * time.Duration(i))
		skipLogin = true
	}

	tests := models.MapTestDataRawData(testsRaw)

	log.Println("Writting Tests to Database")
	newTests := WriteTests(tests, db)
	log.Println("Triggering Notifications")
	triggerTestsNotifications(newTests, tests)
}