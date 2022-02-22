package models

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"gorm.io/gorm"
	"net"
)

type Hardware struct {
	gorm.Model

	ContractID uint `form:"contract_id" json:"contract_id" gorm:"not null"`
	Contract   *Contract

	Cpu  string `form:"cpu" json:"cpu"`
	Disk string `form:"disk" json:"disk"`
	Host string `form:"host" json:"host"`
	Net  string `form:"net" json:"net"`
}

type HardwareInfo struct {
	Cpu  []cpu.InfoStat  `json:"cpu"`
	Disk []string        `json:"disk"`
	Host host.InfoStat   `json:"host"`
	Net  []net.Interface `json:"net"`
}
