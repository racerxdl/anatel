package main

import (
	"github.com/anaskhan96/soup"
	"github.com/tebeka/selenium"
	"log"
	"time"
)

func consultaIndicativoArray(indicativos []string, driver selenium.WebDriver) (result []map[string]string) {
	result = make([]map[string]string, 0)

	defer func() {
		if r := recover(); r != nil {
			// Recover and return partial
			log.Println("Error:", r)
		}
	}()

	for i := 0; i < len(indicativos); i++ {
		log.Println("Fetching data for", indicativos[i])
		result = append(result, consultaIndicativo(indicativos[i], driver)...)
		if i%64 == 63 {
			log.Println("Waiting 5 seconds before next batch")
			time.Sleep(time.Second * 5)
		}
	}
	return result
}

func consultaIndicativo(indicativo string, webDriver selenium.WebDriver) []map[string]string {
	webDriver.SetImplicitWaitTimeout(2 * time.Second)
	webDriver.Get(anatelConsultaIndicativo)
	certs := make([]map[string]string, 0)

	//log.Println("Getting pIndicativo")
	elem, err := webDriver.FindElement(selenium.ByID, "pIndicativo")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(indicativo)

	elem, err = webDriver.FindElement(selenium.ByID, "botaoFlatConfirmar")
	if err != nil {
		panic(err)
	}

	elem.Click()

	SeleniumWait(webDriver, 100*time.Millisecond)

	content, err := webDriver.PageSource()

	doc := soup.HTMLParse(content)

	tables := doc.FindAll("table", "class", "Tabela")

	if len(tables) == 0 {
		panic("Cannot find table with class Tabela")
	}

	for n := 0; n < len(tables); n++ {
		table := tables[n]

		z := table.Find("td", "class", "CampoCentro")

		if z.Error != nil {
			//log.Println("Not find child with CampoCentro")
			continue
		}

		headers := make([]string, 0)

		heads := table.FindAll("th")

		for i := 0; i < len(heads); i++ {
			headers = append(headers, CleanString(heads[i].Text()))
		}

		trs := table.FindAll("tr")

		for i := 0; i < len(trs); i++ {
			d := make(map[string]string)
			tr := trs[i]

			ths := tr.FindAll("th")

			if len(ths) > 0 {
				continue
			}

			tds := tr.FindAll("td")

			for z := 0; z < len(tds); z++ {
				td := tds[z]
				d[headers[z]] = CleanString(td.Text())
			}

			d["__message"] = "OK"
			d["__req"] = indicativo

			certs = append(certs, d)
		}

		if len(trs) == 0 {
			d := make(map[string]string)
			d["__message"] = "NOT_FOUND"
			d["__req"] = indicativo
			certs = append(certs, d)
		}
	}

	return certs
}
