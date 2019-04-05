package main

import (
	"github.com/anaskhan96/soup"
	"github.com/tebeka/selenium"
	"strconv"
	"time"
)

func consultaAgenda(username, password, uf string, startTime time.Time, skipLogin bool, webDriver selenium.WebDriver) []map[string]string {
	webDriver.SetImplicitWaitTimeout(10 * time.Second)
	certs := make([]map[string]string, 0)

	if !skipLogin {
		webDriver.Get(anatelSEC)

		//log.Println("Getting SISLogin")
		elem, err := webDriver.FindElement(selenium.ByID, "SISLogin")
		if err != nil {
			panic(err)
		}

		elem.SendKeys(username)

		//log.Println("Getting SISSenha")
		elem, err = webDriver.FindElement(selenium.ByID, "SISSenha")
		if err != nil {
			panic(err)
		}

		elem.SendKeys(password)

		elem, err = webDriver.FindElement(selenium.ByID, "botaoFlatEntrar")
		if err != nil {
			panic(err)
		}

		elem.Click()

		SeleniumWait(webDriver, 2000*time.Millisecond)
	}

	webDriver.Get(anatelConsultaAgenda)

	SeleniumWait(webDriver, 2000*time.Millisecond)

	endTime := startTime.Add(time.Hour * 24 * 30)

	startDateTime := startTime.Format("02012006")
	endDateTime := endTime.Format("02012006")

	elem, err := webDriver.FindElement(selenium.ByID, "DataInicial")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(startDateTime)

	elem, err = webDriver.FindElement(selenium.ByID, "DataFinal")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(endDateTime)

	elem, err = webDriver.FindElement(selenium.ByID, "UF")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(uf)

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
				d[headers[z]] = CleanString(GetStringNested(&td))
			}

			d["UF"] = uf
			d["__message"] = "OK"
			d["__startDate"] = strconv.FormatInt(startTime.Unix(), 10)
			d["__endDate"] = strconv.FormatInt(startTime.Unix(), 10)

			certs = append(certs, d)
		}
	}

	return certs
}
