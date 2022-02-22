package service

import (
	"licenseGenerator/dao"
	"licenseGenerator/models"
)

//func GetEncryptProductList(license models.License) ([]models.Product, error) {
//	productList, err := GetContractProductList(license.ContractID)
//	if err != nil {
//		return nil, err
//	}
//	err = dao.Db.First(&license.Key, license.KeyID).Error
//	if err != nil {
//		return nil, err
//	}
//	for i, product := range productList {
//		if product.ProductFuncList != nil {
//			productFuncListBytes, _ := json.Marshal(product.ProductFuncList)
//			publicKey, err := DecodePublicKeyString(license.Key.PublicKey)
//			if err != nil {
//				return nil, err
//			}
//			productList[i].ProductEncryptFuncListString, err = Encrypt(publicKey, productFuncListBytes)
//			if err != nil {
//				return nil, err
//			}
//		}
//	}
//	return productList, nil
//}

func GetContractProductList(contractID uint) ([]models.Product, error) {
	var productList []models.Product
	var bindProductList []models.BindProduct
	err := dao.Db.Where("contract_id = ?", contractID).Find(&bindProductList).Error
	if err != nil {
		return nil, err
	}
	for _, bindProduct := range bindProductList {
		var product models.Product
		err = dao.Db.First(&product, bindProduct.ProductID).Error
		if err != nil {
			return nil, err
		}
		var bindFuncList []models.BindFunc
		err = dao.Db.Where("bind_product_id = ?", bindProduct.ID).Find(&bindFuncList).Error
		if err != nil {
			return nil, err
		}
		if bindProductList != nil {
			for _, bindFunc := range bindFuncList {
				var function models.Func
				err = dao.Db.First(&function, bindFunc.FuncID).Error
				if err != nil {
					return nil, err
				}
				product.ProductFuncList = append(product.ProductFuncList, function)
			}
		}
		productList = append(productList, product)
	}
	return productList, nil
}

func GetContractHardwareList(contractID uint) ([]models.Hardware, error) {
	var hardwareList []models.Hardware
	err := dao.Db.Where("contract_id = ?", contractID).Find(&hardwareList).Error
	if err != nil {
		return nil, err
	}
	return hardwareList, nil
}

func SyncCrm(contract *models.Contract) error {
	var crm models.ForPaasContractView
	err := dao.CrmDb.Where("contract_code = ?", contract.ContractCode).First(&crm).Error
	if err != nil {
		return err
	}
	contract.ContractStatus = crm.ContractStatus
	contract.ContractCustomerName = crm.ClientName
	contract.ContractSellerName = crm.ResponbilitierName
	contract.ContractSellerID = crm.ResponbilitierID
	return nil
}
