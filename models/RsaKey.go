package models

import (
	"gorm.io/gorm"
)

type RsaKey struct {
	gorm.Model

	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}
