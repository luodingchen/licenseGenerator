package models

import "gorm.io/gorm"

type BindProduct struct {
	gorm.Model

	ContractID uint `json:"contract_id" gorm:"not null"`
	Contract   *Contract

	ProductID uint `json:"product_id" gorm:"not null"`
	Product   *Product

	BindFuncList []BindFunc `json:"bind_func_list" gorm:"-"`
}
