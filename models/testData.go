package models

import (
	"time"
	"github.com/jinzhu/gorm"
	"strings"
	"strconv"
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

type TestData struct {
	gorm.Model
	UID string
	Hash string
	Certificates string
	TestDate time.Time
	InscriptionEndDate time.Time
	ActiveInscriptions string
	CanceledInscriptions string
	StartTime string
	Address string
	MorseTest bool
	ComputerTest bool
	Accountable string
	Status string
	ContactPhone string
	MaxVacancies string
	AvailableVacancies string
	Region string
}

func (td *TestData) GenerateUID() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strconv.FormatInt(td.TestDate.UnixNano(), 10)))
	hasher.Write([]byte(td.Certificates))
	hasher.Write([]byte(td.Region))
	result = hasher.Sum(result)

	td.UID = strings.ToUpper(hex.EncodeToString(result))
}

func (td *TestData) GenerateHash() {
	var result []byte

	hasher := sha3.New512()
	hasher.Write([]byte(strconv.FormatInt(td.TestDate.UnixNano(), 10)))
	hasher.Write([]byte(td.Certificates))
	hasher.Write([]byte(strconv.FormatInt(td.InscriptionEndDate.Unix(), 10)))
	hasher.Write([]byte(td.ActiveInscriptions))
	hasher.Write([]byte(td.CanceledInscriptions))
	hasher.Write([]byte(strconv.FormatBool(td.MorseTest)))
	hasher.Write([]byte(strconv.FormatBool(td.ComputerTest)))
	hasher.Write([]byte(td.Accountable))
	hasher.Write([]byte(td.AvailableVacancies))
	hasher.Write([]byte(td.Region))
	result = hasher.Sum(result)

	td.Hash = strings.ToUpper(hex.EncodeToString(result))
}