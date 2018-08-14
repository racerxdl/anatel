package main

import (
	"log"
	"github.com/racerxdl/anatel/models"
	"github.com/tebeka/selenium"
	"github.com/jinzhu/gorm"
	"time"
)

func UpdateStationsFlow(username, password, uf string, db *gorm.DB, driver selenium.WebDriver) {
	log.Println("Fetching callsigns")
	callsignsRaw := consultaIndicativos(username, password, uf, ClassC, driver)

	log.Println("Santizing Callsigns")
	callSigns := models.MapCallsignRawData(callsignsRaw[:11])

	log.Println("Writting Callsigns")

	newCallsigns := WriteCallSigns(callSigns, db)

	log.Println("Fetching Extended Data")

	extDataRaw := consultaIndicativoArray(models.CallSignArrayToString(callSigns), driver)
	extData := models.MapStationRawData(extDataRaw)

	log.Println("Writting Station Data")
	newStations := WriteStationData(extData, db)

	triggerStationCallSignsNotifications(newCallsigns, newStations, callSigns, extData, db)
}

const testCheckMonths = 6

func GetNextTests(username, password, uf string, db *gorm.DB, driver selenium.WebDriver) {
	startTime := time.Now().Add(-time.Hour * 24 * 30)

	startTime = startTime.Add(-time.Hour * 24 * time.Duration(startTime.Day() - 1)) // Reset for the first day of month

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