package eventmanager

import "github.com/racerxdl/anatel/models"

const EvOnNewCallsign = "newCallsignEvent"

type NewCallsignEventData struct {
	Owner    string
	CallSign models.CallSign
	Stations []models.StationData
}
