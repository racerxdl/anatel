package main

import (
	"log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
	"github.com/racerxdl/anatel/eventmanager"
)


var newCallsign = make(chan interface{})
var newStation = make(chan interface{})
var newTests = make(chan interface{})

func OnNewCallsign(data eventmanager.NewCallsignEventData) {
	log.Printf("New Callsign %s with %d stations!\n", data.CallSign.Callsign, len(data.Stations))
	for i := 0; i < len(data.Stations); i++ {
		st := data.Stations[i]
		log.Printf("    Station: %s - %s from %s\n", st.Callsign, st.StationType, st.Owner)
	}
}

func OnNewStation(data eventmanager.NewStationEventData) {
	log.Printf("New Station  %s - %s from %s\n", data.Callsign, data.StationType, data.Owner)
}

func OnNewTestDate(data eventmanager.NewTestDateEventData) {
	log.Printf("New test in %s at %s for %s\n", data.TestDate.String(), data.StartTime, data.Certificates)
}

func main() {
	hubUrl := "http://localhost:4444/wd/hub"

	var webDriver selenium.WebDriver
	var err error

	eventManager.AddHandler(eventmanager.EvOnNewCallsign, newCallsign)
	eventManager.AddHandler(eventmanager.EvOnNewStation, newStation)
	eventManager.AddHandler(eventmanager.EvOnNewTestDate, newTests)

	go func() {
		log.Println("Starting Handler loop")
		for {
			select {
			case msg := <-newCallsign:
				OnNewCallsign(msg.(eventmanager.NewCallsignEventData))
			case msg := <-newStation:
				OnNewStation(msg.(eventmanager.NewStationEventData))
			case msg := <-newTests:
				OnNewTestDate(msg.(eventmanager.NewTestDateEventData))
			}
		}
		log.Println("Ending Handler loop")
	}()

	db := Initialize()

	defer db.Close()

	caps := selenium.Capabilities(map[string]interface{}{
		"browserName": "chrome",
		"enableVideo": true,
		"screenResolution": "1280x1024x24",
	})

	caps.AddChrome(chrome.Capabilities{})

	log.Println("Initializing Remote")
	if webDriver, err = selenium.NewRemote(caps, hubUrl); err != nil {
		panic(err)
	}

	defer webDriver.Quit()

	window, _ := webDriver.CurrentWindowHandle()

	log.Println("Resizing Window")
	webDriver.SetImplicitWaitTimeout(5 * time.Second)
	webDriver.ResizeWindow(window, 1280, 1024)

	// UpdateStationsFlow(db, webDriver)

	//z := consultaAgenda("-----", "----", "SP", webDriver)

	GetNextTests("----", "----",  "SP", db, webDriver)

	//z := consultaCertificado("-----", webDriver)
	//z :=  consultaIndicativos("----", "----", "SP", ClassC, webDriver)
	//z := consultaIndicativoArray([]string {"PY2KSC", "PY2KJP"}, webDriver)
	//v, _ := json.MarshalIndent(z, "", "    ")
	//log.Println(string(v))
}