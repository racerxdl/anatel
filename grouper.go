package main

import (
	"encoding/json"
	"github.com/tebeka/selenium"
	"log"
)

func RetrieveCallSign(username, password, uf string, driver selenium.WebDriver) {
	log.Println("Fetching list of callSigns")
	callsignsObj := consultaIndicativos(username, password, uf, ClassC, driver)
	log.Println("Found", len(callsignsObj), "callsigns")
	callsigns := Map2Str(callsignsObj, func(f map[string]string) string {
		return f["indicativo"]
	})

	vData := consultaIndicativoArray(callsigns, driver)

	v, _ := json.MarshalIndent(vData, "", "    ")
	log.Println(string(v))
}
