package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"strings"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

type RepeaterStationData struct {
	gorm.Model
	UID string
	RXFrequency uint64
	TXFrequency uint64
	Callsign string
	StationNumber string
	City string
	Region string
	StationType string
	FirstSaw time.Time
}

func (cs *RepeaterStationData) GenerateUID() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strings.ToLower(cs.Callsign)))
	hasher.Write([]byte(strings.ToLower(cs.StationNumber)))
	hasher.Write([]byte(strings.ToLower(cs.StationType)))

	result = hasher.Sum(result)

	cs.UID = strings.ToUpper(hex.EncodeToString(result))
}