package eventmanager

import "github.com/racerxdl/anatel/models"

const EvOnNewStation = "newStationEvent"

type NewStationEventData struct {
	models.StationData
}
