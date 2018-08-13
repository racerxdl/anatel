package main

import (
	"log"
	"github.com/tebeka/selenium"
)

func consultaCertificado(cpfCnpj string, webDriver selenium.WebDriver) []map[string]string {
	webDriver.Get(anatelConsultaCertificadoURL)
	certs := make([]map[string]string, 0)

	log.Println("Getting pNumCnpjCpf")
	elem, err := webDriver.FindElement(selenium.ByID, "pNumCnpjCpf")
	if err != nil {
		panic(err)
	}

	elem.SendKeys(cpfCnpj)

	elem, err = webDriver.FindElement(selenium.ByID, "botaoFlatConfirmar")
	if err != nil {
		panic(err)
	}

	elem.Click()

	elem, err = webDriver.FindElement(selenium.ByID, "divestacao")
	if err != nil {
		panic(err)
	}

	elems, err := elem.FindElements(selenium.ByClassName, "Tabela")

	if err != nil {
		panic(err)
	}

	for n := 0; n < len(elems); n++ {
		table := elems[n]

		z, err := table.FindElements(selenium.ByClassName, "CampoCentro")

		if len(z) == 0 || err != nil {
			continue
		}

		elems, err = table.FindElements(selenium.ByTagName, "tr")

		if err != nil {
			panic(err)
		}

		if len(elems) < 1 {
			log.Println("Buggy")
			continue
		}

		trHeader := elems[0]
		elems = elems[1:]

		thNames := make([]string, 0)

		ths, err := trHeader.FindElements(selenium.ByTagName, "th")

		if err != nil {
			panic(err)
		}

		for i := 0; i < len(ths); i++ {
			s, err := ths[i].Text()
			if err != nil {
				panic(err)
			}
			thNames = append(thNames, s)
		}

		for i := 0; i < len(elems); i++ {
			certData := make(map[string]string)
			v := elems[i]
			c, err := v.FindElements(selenium.ByClassName, "CampoCentro")

			if err != nil {
				panic(err)
			}

			for z := 0; z < len(c); z++ {
				certData[thNames[z]], _ = c[z].Text()
			}

			certs = append(certs, certData)
		}

		break
	}

	return certs
}