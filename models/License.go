package models

import (
	"gorm.io/gorm"
)

type License struct {
	gorm.Model `yaml:"-"`

	ContractID uint      `json:"contract_id" gorm:"not null" yaml:"-"`
	Contract   *Contract `yaml:"-"`

	KeyID uint    `json:"key_id" gorm:"not null"`
	Key   *RsaKey `yaml:"-"`

	LicenseType         uint   `json:"license_type" yaml:"-"`   // 0试用、1正式99年 试用版期限最高30天
	LicenseStatus       uint   `json:"license_status" yaml:"-"` // 0:待授权 1:已授权 2:已签名 3:已拒绝
	LicenseDuration     int    `json:"license_duration" yaml:"-"`
	LicenseStartTime    string `json:"license_start_time" yaml:"-"`
	LicenseRejectedTime string `json:"license_rejected_time" yaml:"-"`
	LicenseDeadline     string `json:"license_deadline" yaml:"-"`
	LicenseDescription  string `json:"license_description" yaml:"-"`
	Config              Config `json:"config" yaml:"-" gorm:"-"`
	ConfigJson          string `json:"config_json" yaml:"config_json"`
	LicensePublicKey    string `json:"license_public_key" yaml:"license_public_key"`
	LicenseSignature    string `json:"license_signature" yaml:"license_signature"`
}

type Config struct {
	StartTime    string     `json:"start_time"`
	Deadline     string     `json:"deadline"`
	ProductList  []Product  `json:"product_list"`
	HardwareList []Hardware `json:"hardware_list"`
}
