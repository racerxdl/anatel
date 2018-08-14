package main

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
	"github.com/racerxdl/anatel/eventmanager"
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
	anatelConsultaIndicativo = "https://sistemas.anatel.gov.br/easp/Novo/ConsultaIndicativo/Tela.asp?SISQSmodulo=11265"
	anatelSEC = "https://sistemas.anatel.gov.br/SEC/"
	anatelConsultaAgenda = "https://sistemas.anatel.gov.br/SEC/Agenda/Tela.asp?OP=c&SISQSmodulo=5819"
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

func Map2Str(vs []map[string]string, f func(map[string]string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

var eventManager = eventmanager.New()
