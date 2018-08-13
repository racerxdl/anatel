package main

import (
	"log"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"time"
)

func main() {
	hubUrl := "http://localhost:4444/wd/hub"

	var webDriver selenium.WebDriver
	var err error

	caps := selenium.Capabilities(map[string]interface{}{
		"browserName": "chrome",
		"enableVideo": true,
		"screenResolution": "1024x768x24",
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
	webDriver.ResizeWindow(window, 1024, 768)

	//z := consultaCertificado("---", webDriver)
	//z :=  consultaIndicativos("----", "----", "SP", ClassC, webDriver)

	//v, _ := json.MarshalIndent(z, "", "    ")
	//log.Println(string(v))
}