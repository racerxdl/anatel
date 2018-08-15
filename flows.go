package main

import (
	"log"
	"github.com/racerxdl/anatel/models"
	"github.com/tebeka/selenium"
	"github.com/jinzhu/gorm"
	"time"
	"fmt"
)

func UpdateStationsFlow(username, password, uf string, db *gorm.DB, driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered", r)
		}
	}()
	driver.DeleteAllCookies()
	log.Println("Fetching callsigns")
	callsignsRaw := consultaIndicativos(username, password, uf, ClassC, driver)

	log.Println("Santizing Callsigns")
	callSigns := models.MapCallsignRawData(callsignsRaw)

	nc := make(chan []int)

	go func() {
		log.Println("Writting Callsigns")
		newCallsigns := WriteCallSigns(callSigns, db)
		nc <- newCallsigns
	}()

	log.Println("Fetching Extended Data")

	extDataRaw := consultaIndicativoArray(models.CallSignArrayToString(callSigns), driver)
	extData := models.MapStationRawData(extDataRaw)

	ns := make(chan []int)

	go func() {
		log.Println("Writting Station Data")
		newStations := WriteStationData(extData, db)
		ns <- newStations
	}()

	log.Println("Waiting for Station / Callsign data to be written")

	newCallsigns := <- nc
	newStations := <- ns

	log.Println("Triggering Notifications")

	triggerStationCallSignsNotifications(newCallsigns, newStations, callSigns, extData, db)
}

const testCheckMonths = 6

func GetNextTests(username, password, uf string, db *gorm.DB, driver selenium.WebDriver) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
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