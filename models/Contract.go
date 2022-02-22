package models

import (
	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model

	ContractCode                string `json:"contract_code" gorm:"not null,unique"` // crm获取
	ContractCustomerName        string `json:"contract_customer_name"`
	ContractSellerName          string `json:"contract_seller_name"`
	ContractSellerID            uint   `json:"contract_seller_id"`
	ContractStatus              uint   `json:"contract_status"`               //0:已关闭  从crm获取，只对关闭前的合同进行授权
	ContractTrailLicenseStatus  bool   `json:"contract_trail_license_status"` // 该合同的试用状态，true表示已生成过试用证书
	ContractFormalLicenseStatus bool   `json:"contract_formal_license_status"`
	ContractDescription         string `json:"contract_description"`
	//ContractFuncList     []Func     `json:"contract_func_list" gorm:"-"`
	ContractProductList  []Product  `json:"contract_product_list" gorm:"-"`
	ContractHardwareList []Hardware `json:"contract_hardware_list" gorm:"-"`
}

type ForPaasContractView struct {
	ContractCode       string
	ContractName       string
	ContractStatus     uint
	ClientName         string
	ResponbilitierName string
	ResponbilitierID   uint
}

func (ForPaasContractView) TableName() string {
	return "for_paas_contract_view"
}
