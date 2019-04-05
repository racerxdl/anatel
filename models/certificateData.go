package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type CertificateData struct {
	gorm.Model
	TaxId          string
	Class          string
	Certificate    string
	InclusionDate  time.Time
	ExpirationDate time.Time
	Name           string
	Status         string
}
