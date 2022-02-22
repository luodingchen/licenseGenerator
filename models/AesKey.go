package models

import "gorm.io/gorm"

type AesKey struct {
	gorm.Model
	AesKeyString string `json:"aesKeyString"`
}
