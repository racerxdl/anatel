package models

import "time"

const (
	ownerTag = "Autorizado(a)"
	certificateTag = "Certificado"
	certificatesTag = "Certificados"
	testDataTag = "Data da prova"
	inscriptionEndDateTag = "Encerramento das Inscrições"
	activeInscriptionsTag = "Inscrições ativas"
	canceledInscriptionsTag = "Inscrições canceladas"
	startTimeTag = "Início"
	addressTag = "Local"
	morseTestTag = "Prova de Morse"
	computerTestTag = "Prova no Computador"
	accountableTag = "Responsável"
	testStatusTag = "Situação da Agenda"
	contactPhoneTag = "Telefone de Contato"
	maxVacanciesTag = "Total de vagas"
	availableVacancies = "Vagas disponíveis"
	callsignTag = "Indicativo"
	cityTag = "Município"
	inclusionDateTag = "Data de Inclusão"
	expirationDateTag = "Data de Validate"
	serviceTag = "Serviço"
	stationTypeTag = "Tipo Estação"
	regionTag = "UF"
	taxIdTag = "CPF"
	classTag = "Categoria/Classe"
	statusTag = "Situação"
	nameTag = "Nome"
)

const (
	dateFormat = "02/01/2006"
)

func IndexOfString(item string, arr []string) int {
	for idx, n := range arr {
		if n == item {
			return idx
		}
	}

	return -1
}

func Raw2TestData(raw map[string]string) TestData {
	testDate, _ := time.Parse(dateFormat, raw[testDataTag])
	inscriptionEndDate, _ := time.Parse(dateFormat, raw[inscriptionEndDateTag])

	td := TestData{
		Certificates: raw[certificatesTag],
		TestDate: testDate,
		InscriptionEndDate: inscriptionEndDate,
		ActiveInscriptions: raw[activeInscriptionsTag],
		CanceledInscriptions: raw[canceledInscriptionsTag],
		StartTime: raw[startTimeTag],
		Address: raw[addressTag],
		MorseTest: raw[morseTestTag] == "Sim",
		ComputerTest: raw[computerTestTag] == "Sim",
		Accountable: raw[accountableTag],
		Status: raw[testStatusTag],
		ContactPhone: raw[contactPhoneTag],
		MaxVacancies: raw[maxVacanciesTag],
		AvailableVacancies: raw[availableVacancies],
		Region: raw[regionTag],
	}

	td.GenerateUID()
	td.GenerateHash()

	return td
}

func RemoveDuplicateTestData(data []TestData) []TestData {
	addedIds := make([]string, 0)
	filteredData := make([]TestData, 0)

	for i := 0; i < len(data); i++ {
		c := data[i]
		if IndexOfString(c.UID, addedIds) == -1 {
			addedIds = append(addedIds, c.UID)
			filteredData = append(filteredData, c)
		}
	}

	return filteredData
}

func MapTestDataRawData(raw []map[string]string) []TestData {
	var cls = make([]TestData, 0)

	for i := 0; i < len(raw); i++ {
		cls = append(cls, Raw2TestData(raw[i]))
	}

	return RemoveDuplicateTestData(cls)
}

func Raw2StationData(raw map[string]string) StationData {

	inclusionDate, _ := time.Parse(dateFormat, raw[inclusionDateTag])
	expirationDate, _ := time.Parse(dateFormat, raw[expirationDateTag])

	station := StationData{
		Owner: raw[ownerTag],
		Certificate: raw[certificateTag],
		InclusionDate: inclusionDate,
		ExpirationDate: expirationDate,
		Callsign: raw[callsignTag],
		City: raw[cityTag],
		Service: raw[serviceTag],
		StationType: raw[stationTypeTag],
		Region: raw[regionTag],
		FirstSaw: time.Now(),
	}

	station.GenerateUID()

	return station
}


func MapStationRawData(raw []map[string]string) []StationData {
	var cls = make([]StationData, 0)

	for i := 0; i < len(raw); i++ {
		cls = append(cls, Raw2StationData(raw[i]))
	}

	return cls
}

func Raw2CallsignData(raw map[string]string) CallSign {
	return CallSign{
		Callsign: raw[callsignTag],
		Region: raw[regionTag],
		Class: raw[classTag],
		LastUpdated: time.Now(),
		FirstSaw: time.Now(),
	}
}

func MapCallsignRawData(raw []map[string]string) []CallSign {
	var cls = make([]CallSign, 0)

	for i := 0; i < len(raw); i++ {
		cls = append(cls, Raw2CallsignData(raw[i]))
	}

	return cls
}

func Raw2CertificateData(raw map[string]string) CertificateData {
	inclusionDate, _ := time.Parse(dateFormat, raw[inclusionDateTag])
	expirationDate, _ := time.Parse(dateFormat, raw[expirationDateTag])

	return CertificateData{
		TaxId: raw[taxIdTag],
		Class: raw[classTag],
		Certificate: raw[certificateTag],
		InclusionDate: inclusionDate,
		ExpirationDate: expirationDate,
		Name: raw[nameTag],
		Status: raw[statusTag],
	}
}
