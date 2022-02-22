package models

import "gorm.io/gorm"

type Func struct {
	gorm.Model

	ProductID 		uint  		`form:"product_id" json:"product_id" gorm:"not null"`
	Product   		*Product
// 在做外键关联时，shouldBind要保证子表与父表没有相同名称字段
// 某些情况若要使用string作为外键关联时，gorm的tag中需要指定长度例如char(36)，因为mysql数据库需要通过定长字段来关联表
	FunctionCode  	string 		`form:"function_code" json:"function_code" gorm:"not null"`
	FunctionName  	string 		`form:"function_name" json:"function_name"`    // 可以为中文
	FunctionValue 	string 		`form:"function_value" json:"function_value"`  // 可以为中文
}
