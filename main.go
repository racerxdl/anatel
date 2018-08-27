package main

import (
	"log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
	"github.com/racerxdl/anatel/eventmanager"
	"fmt"
	"os"
	"strings"
	"strconv"
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
	telBot.SendMessage(
		fmt.Sprintf("Nova prova em *%s* para _%s_ ! Ela começa as *%s* do dia *%s* e será em *%s*.\nTotal de Vagas: *%s*\nEncerramento das Inscrições: *%s*",
			data.Region,
			data.Certificates,
			data.StartTime,
			data.TestDate.Format("02/01/2006"),
			data.Address,
			data.MaxVacancies,
			data.InscriptionEndDate.Format("02/01/2006"),
		),
	)
}

func CreateWebDriver() selenium.WebDriver {
	var webDriver selenium.WebDriver
	var enableVideoStr = os.Getenv("ENABLE_VIDEO")
	var enableVideo = false

	if enableVideoStr != "" {
		var err error
		enableVideo, err = strconv.ParseBool(enableVideoStr)
		if err != nil {
			panic(fmt.Sprintf("Error parsing ENABLE_VIDEO: %s", err))
		}
	}

	hubUrl := os.Getenv("SELENOID_URL")

	var err error
	caps := selenium.Capabilities(map[string]interface{}{
		"browserName": "chrome",
		"enableVideo": enableVideo,
		"screenResolution": "1280x1024x24",
	})

	caps.AddChrome(chrome.Capabilities{})

	log.Println("Initializing Remote")
	if webDriver, err = selenium.NewRemote(caps, hubUrl); err != nil {
		panic(err)
	}

	return webDriver
}

func main() {
	mode := os.Getenv("MODE")
	database = Initialize()
	defer database.Close()

	if mode == "graphql" {
		graphqlServer()
		return
	}

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


	webDriver := CreateWebDriver()

	defer func() {
		if webDriver != nil {
			webDriver.Quit()
		}
	}()

	window, _ := webDriver.CurrentWindowHandle()

	log.Println("Resizing Window")
	webDriver.SetImplicitWaitTimeout(5 * time.Second)
	webDriver.ResizeWindow(window, 1280, 1024)

	state := os.Getenv("STATE")

	checkstateRaw := strings.Split(state, ",")
	checkstates := make([]string, 0)

	for i := 0; i < len(checkstateRaw); i++ {
		s := CleanString(checkstateRaw[i])

		if len(s) == 0 || IndexOfString(s, States) == -1 {
			log.Println("Invalid state: ", s)
			continue
		}

		checkstates = append(checkstates, s)
	}

	if len(checkstates) == 0 {
		panic("No valid states provided.")
	}

	log.Println("Starting checks for state", state, "with mode", mode)

	if mode == "tests" || mode == "all" {
		log.Println("Checking next 6 month tests")
		// region Get Tests
		for i := 0; i < len(checkstates); i++ {
			state := checkstates[i]
			log.Println("Checking tests for", state)
			GetNextTests(os.Getenv("ANATEL_USERNAME"), os.Getenv("ANATEL_PASSWORD"), state, database, webDriver)
			webDriver.DeleteAllCookies() // Force login again
		}
	}

	// endregion

	var classes = []string {ClassC, ClassB, ClassA}

	if mode == "callsign" || mode == "all" {
		// region Update Callsigns
		log.Println("Checking callsigns")
		for i := 0; i < len(checkstates); i++ {
			state := checkstates[i]
			for z := 0; z < len(classes); z++ {
				class := classes[z]
				log.Println("Checking callsigns for", state, class)
				// region Class C
				UpdateCallSigns(os.Getenv("ANATEL_USERNAME"), os.Getenv("ANATEL_PASSWORD"), state, class, database, webDriver)
				webDriver.DeleteAllCookies()
				webDriver.Close()
				webDriver = CreateWebDriver()
				UpdateStationsFlow(state, database, webDriver)
				webDriver.DeleteAllCookies()
				webDriver.Close()
				webDriver = CreateWebDriver()
				// endregion
			}
		}

		// endregion
	}

	if mode == "repeater" || mode == "all" {
		// region Update Repeaters
		log.Println("Checking repeaters")
		for i := 0; i < len(checkstates); i++ {
			state := checkstates[i]
			log.Println("Checking repeaters for", state)
			UpdateRepeaterStationsFlow(os.Getenv("ANATEL_USERNAME"), os.Getenv("ANATEL_PASSWORD"), state, database, webDriver)
			webDriver.DeleteAllCookies()
			webDriver.Close()
			webDriver = CreateWebDriver()
		}
		// endregion
	}

	log.Println("Finished all tasks. Waiting for notifications")
	time.Sleep(time.Second * 60)
	log.Println("Closing")
}