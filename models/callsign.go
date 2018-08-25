package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type CallSign struct {
	gorm.Model
	Callsign string `gorm:"primary_key"`
	FirstSaw time.Time
	LastUpdated time.Time
	Class string
	Region string
	Stations []StationData `gorm:"foreignkey:Callsign;association_foreignkey:Callsign"`
	Repeaters []RepeaterStationData `gorm:"foreignkey:Callsign;association_foreignkey:Callsign"`
}

func CallSignArrayToString(data []CallSign) []string {
	arr := make([]string, 0)
	for i := 0; i < len(data); i++ {
		arr = append(arr, data[i].Callsign)
	}
	return arr
}