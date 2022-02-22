package models

import "gorm.io/gorm"

type BindFunc struct {
	gorm.Model

	BindProductID uint `json:"bind_product_id" gorm:"not null"`
	BindProduct   *BindProduct

	FuncID uint `json:"func_id" gorm:"not null"`
	Func   *Func
}
