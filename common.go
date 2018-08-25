package main

import (
	"github.com/tebeka/selenium"
	"strings"
	"time"
	"github.com/racerxdl/anatel/eventmanager"
	"github.com/racerxdl/anatel/telegram"
)

const (
	ClassC = "Classe C"
	ClassB = "Classe B"
	ClassA = "Classe A"
)

var States = []string {
	"AC", "AL", "AM", "AP", "BA", "CE", "DF",
	"ES", "GO", "MA", "MG", "MS", "MT", "PA",
	"PB", "PE", "PI", "PR", "RJ", "RN", "RO",
	"RR", "RS", "SC", "SE", "SP", "TO",
}

const callSignUpdateTimeout = time.Hour * 24 * 7 // 1 week
const segmentLength = 100

const (
	anatelSCRAURL = "https://sistemas.anatel.gov.br/SCRA/"
	anatelSCRAIndicativo = "https://sistemas.anatel.gov.br/SCRA/ConsultaIndicativoVagoOcupado/tela.asp?SISQSmodulo=18082"
	anatelSCRARepetidora = "https://sistemas.anatel.gov.br/SCRA/Relatorio/Repetidora.asp?SISQSmodulo=16442"
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
var telBot = telegram.New()
