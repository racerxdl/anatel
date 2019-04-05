package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/tebeka/selenium"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"
)

func consultaIndicativos(username, password, uf, class string, webDriver selenium.WebDriver) []map[string]string {
	webDriver.SetImplicitWaitTimeout(10 * time.Second)
	webDriver.Get(anatelSCRAURL)
	certs := make([]map[string]string, 0)

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

	//log.Println("Acessing Indicativos")
	err = webDriver.Get(anatelSCRAIndicativo)

	if err != nil {
		panic(err)
	}

	elem, err = webDriver.FindElement(selenium.ByID, "SiglaUF")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(uf)

	elem, err = webDriver.FindElement(selenium.ByID, "idtTipoIndicativo")
	if err != nil {
		panic(err)
	}

	elem.SendKeys("Efetivo")

	elem, err = webDriver.FindElement(selenium.ByID, "CodCategoria")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(class)

	elem, err = webDriver.FindElement(selenium.ByID, "IndOcupado1")
	if err != nil {
		panic(err)
	}

	elem.Click()

	elem, err = webDriver.FindElement(selenium.ByID, "botaoFlatConfirmar")
	if err != nil {
		panic(err)
	}

	elem.Click()

	elems, err := webDriver.FindElements(selenium.ByClassName, "SubTituloEsquerda")
	if err != nil {
		panic(err)
	}

	var regElem selenium.WebElement

	for i := 0; i < len(elems); i++ {
		s, err := elems[i].Text()
		if err != nil {
			panic(err)
		}
		if matched, err := regexp.Match("Registro", []byte(s)); matched {
			if err != nil {
				panic(err)
			}

			regElem = elems[i]
			break
		}
	}

	if regElem == nil {
		panic("Cannot find registro")
	}

	elems, err = regElem.FindElements(selenium.ByClassName, "TextoAzul2")

	higherNumber := 0

	for i := 0; i < len(elems); i++ {
		s, _ := elems[i].Text()
		v, err := strconv.ParseInt(s, 10, 32)

		if err != nil {
			panic(err)
		}

		higherNumber = int(math.Max(float64(higherNumber), float64(v)))
	}

	//log.Println("Number of registros:", higherNumber)

	elem, err = webDriver.FindElement(selenium.ByID, "NumReg")

	if err != nil {
		panic(err)
	}

	elem.SendKeys(fmt.Sprintf("%d", higherNumber))

	_, err = webDriver.ExecuteScript("AlteraNumReg();", nil)

	if err != nil {
		panic(err)
	}

	SeleniumWait(webDriver, 100*time.Millisecond)

	_, err = webDriver.ExecuteScript("AlteraNumReg();", nil)

	if err != nil {
		panic(err)
	}

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

		tds := table.FindAll("td", "class", "CampoCentro")

		log.Println("Enumerating indicativos")
		for i := 0; i < len(tds); i++ {
			td := tds[i]

			s := td.Text()

			certs = append(certs, map[string]string{
				"Indicativo":       CleanString(s),
				"UF":               uf,
				"Categoria/Classe": class,
			})
		}
	}

	return certs
}
