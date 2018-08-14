package eventmanager

import "github.com/racerxdl/anatel/models"

const EvOnNewTestDate = "newTestDateEvent"

type NewTestDateEventData struct {
	models.TestData
}
