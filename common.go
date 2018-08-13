package main

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

const (
	ClassC = "Classe C"
	ClassB = "Classe B"
	ClassA = "Classe A"
)

const (
	anatelSCRAURL = "https://sistemas.anatel.gov.br/SCRA/"
	anatelSCRAIndicativo = "https://sistemas.anatel.gov.br/SCRA/ConsultaIndicativoVagoOcupado/tela.asp?SISQSmodulo=18082"
	anatelConsultaCertificadoURL = "https://sistemas.anatel.gov.br/easp/Novo/ConsultaCertificado/Tela.asp?SISQSmodulo=19176"
)

func SeleniumWait(webDriver selenium.WebDriver, interval time.Duration) {
	webDriver.WaitWithTimeout(Nothing, interval)
}

func Nothing (_ selenium.WebDriver) (bool, error) {
	return false, nil
}

func CleanString(s string) string {
	return strings.Replace(strings.Replace(s, "\t", "", -1), "\n", "", -1)
}