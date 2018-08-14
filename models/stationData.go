package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"golang.org/x/crypto/sha3"
	"strconv"
	"strings"
	"encoding/hex"
)

type StationData struct {
	gorm.Model
	UID string
	Owner string
	Certificate string
	InclusionDate time.Time
	ExpirationDate time.Time
	Callsign string
	City string
	Service string
	StationType string
	Region string
	FirstSaw time.Time
}

func (cs *StationData) GenerateUID() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strings.ToLower(cs.Callsign)))
	hasher.Write([]byte(strconv.FormatInt(cs.InclusionDate.UnixNano(), 10)))
	hasher.Write([]byte(strconv.FormatInt(cs.ExpirationDate.UnixNano(), 10)))
	hasher.Write([]byte(strings.ToLower(cs.StationType)))
	hasher.Write([]byte(strings.ToLower(cs.Owner)))

	result = hasher.Sum(result)

	cs.UID = strings.ToUpper(hex.EncodeToString(result))
}