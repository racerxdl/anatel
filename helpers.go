package main

import (
	"github.com/anaskhan96/soup"
	"github.com/jinzhu/gorm"
	"github.com/racerxdl/anatel/eventmanager"
	"github.com/racerxdl/anatel/models"
)

func IndexOfString(item string, arr []string) int {
	for idx, n := range arr {
		if n == item {
			return idx
		}
	}

	return -1
}

func GetStringNested(n *soup.Root) string {
	str := n.Text()

	c := n.Children()

	for i := 0; i < len(c); i++ {
		v := c[i]
		if v.Error == nil {
			nst := GetStringNested(&v)
			if len(nst) != 0 {
				if len(str) != 0 {
					str += " "
				}
				str += nst
			}
		}
	}

	return str
}

func triggerStationCallSignsNotifications(newCallsigns, newStations []int, callsigns []models.CallSign, stations []models.StationData, db *gorm.DB) {
	callSignsToReport := make([]string, 0)

	for i := 0; i < len(newCallsigns); i++ {
		clsId := newCallsigns[i]
		cls := callsigns[clsId]
		callSignsToReport = append(callSignsToReport, cls.Callsign)
	}

	// Report new callsigns with all stations
	for i := 0; i < len(callSignsToReport); i++ {
		var cls models.CallSign
		cstr := callSignsToReport[i]
		stations := make([]models.StationData, 0)
		db.Model(&models.CallSign{}).Where("callsign = ?", cstr).First(&cls)
		db.Model(&models.StationData{}).Where("callsign = ?", cstr).Find(&stations)

		name := "Desconhecido"

		if len(stations) > 0 {
			name = stations[0].Owner
		}

		eventManager.Emit(eventmanager.EvOnNewCallsign, eventmanager.NewCallsignEventData{
			CallSign: cls,
			Stations: stations,
			Owner:    name,
		})
	}

	// Report new stations for existing callsigns
	for i := 0; i < len(newStations); i++ {
		ssid := newStations[i]
		s := stations[ssid]

		if IndexOfString(s.Callsign, callSignsToReport) == -1 {
			// Not reported in new callsign
			eventManager.Emit(eventmanager.EvOnNewStation, eventmanager.NewStationEventData{
				StationData: s,
			})
		}
	}
}

func triggerStationNotifications(newStations []int, stations []models.StationData, db *gorm.DB) {
	// Report new stations for existing callsigns
	for i := 0; i < len(newStations); i++ {
		ssid := newStations[i]
		s := stations[ssid]

		// Not reported in new callsign
		eventManager.Emit(eventmanager.EvOnNewStation, eventmanager.NewStationEventData{
			StationData: s,
		})
	}
}

func triggerTestsNotifications(newTests []int, tests []models.TestData) {
	for i := 0; i < len(newTests); i++ {
		tdid := newTests[i]
		td := tests[tdid]

		eventManager.Emit(eventmanager.EvOnNewTestDate, eventmanager.NewTestDateEventData{
			TestData: td,
		})
	}
}
