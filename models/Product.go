package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	ProductCode        string `json:"product_code" gorm:"not null"` //字母 -
	ProductName        string `json:"product_name"`
	ProductVersion     string `json:"product_version"` // 1.0 2.0
	ProductStatus      uint   `json:"product_status"`  // 产品关闭时不再售卖 0开启 1关闭
	ProductDescription string `json:"product_description"`
	ProductFuncList    []Func `json:"product_func_list" gorm:"-"`
}
