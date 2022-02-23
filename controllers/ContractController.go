package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"licenseGenerator/service"
)

func PostContract(c *gin.Context) {
	var res = NewResultMsg(c)
	var contract models.Contract
	err := c.ShouldBind(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	if contract.ContractProductList == nil {
		res.Error("Product can not be empty")
		return
	}
	err = service.SyncCrm(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	if contract.ContractStatus == 0 {
		res.Error("contract is closed , can not authorize")
		return
	}
	dao.Db.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		err = tx.Create(&contract).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		for i, product := range contract.ContractProductList {
			var bindProduct models.BindProduct
			bindProduct.ContractID = contract.ID
			bindProduct.ProductID = product.ID

			err = tx.Create(&bindProduct).Error
			if err != nil {
				res.Error(err.Error())
				return err
			}

			err = tx.First(&contract.ContractProductList[i], product.ID).Error
			if err != nil {
				res.Error(err.Error())
				return err
			}

			if contract.ContractProductList[i].ProductFuncList != nil {
				for j, function := range contract.ContractProductList[i].ProductFuncList {
					var bindFunc models.BindFunc
					bindFunc.BindProductID = bindProduct.ID
					bindFunc.FuncID = function.ID

					err = tx.Create(&bindFunc).Error
					if err != nil {
						res.Error(err.Error())
						return err
					}

					err = tx.First(&contract.ContractProductList[i].ProductFuncList[j], function.ID).Error
					if err != nil {
						res.Error(err.Error())
						return err
					}
				}
			}
		}

		if contract.ContractHardwareList != nil {
			var aesKey models.AesKey
			dao.Db.Last(&aesKey)
			for i, hardware := range contract.ContractHardwareList {
				cpuBytes, err := service.AesDecrypt(hardware.Cpu, aesKey.AesKeyString)
				if err != nil {
					res.Error("Decrypt failed , please download latest HardwareMsgGenerator to encrypt your hardware info")
					return err
				}
				diskBytes, err := service.AesDecrypt(hardware.Disk, aesKey.AesKeyString)
				if err != nil {
					res.Error("Decrypt failed , please download latest HardwareMsgGenerator to encrypt your hardware info")
					return err
				}
				hostBytes, err := service.AesDecrypt(hardware.Host, aesKey.AesKeyString)
				if err != nil {
					res.Error("Decrypt failed , please download latest HardwareMsgGenerator to encrypt your hardware info")
					return err
				}
				netBytes, err := service.AesDecrypt(hardware.Net, aesKey.AesKeyString)
				if err != nil {
					res.Error("Decrypt failed , please download latest HardwareMsgGenerator to encrypt your hardware info")
					return err
				}
				contract.ContractHardwareList[i] = models.Hardware{
					ContractID: contract.ID,
					Cpu:        string(cpuBytes),
					Disk:       string(diskBytes),
					Host:       string(hostBytes),
					Net:        string(netBytes),
				}
			}
			err = tx.Create(&contract.ContractHardwareList).Error
			if err != nil {
				res.Error(err.Error())
				return err
			}
		}
		res.Success(PostSuccess, contract)
		return nil
	})
}

func PutContract(c *gin.Context) {
	var res = NewResultMsg(c)
	var contract models.Contract
	err := c.ShouldBind(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = service.SyncCrm(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Updates(&contract).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PutSuccess, contract)
}

func DeleteContract(c *gin.Context) {
	res := NewResultMsg(c)
	var contract models.Contract
	err := c.ShouldBind(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}

	var bindProductList []models.BindProduct
	err = dao.Db.Where("contract_id = ?", contract.ID).Find(&bindProductList).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	dao.Db.Transaction(func(tx *gorm.DB) error {
		for _, bindProduct := range bindProductList {
			err = tx.Where("bind_product_id = ?", bindProduct.ID).Delete(&models.BindFunc{}).Error
			if err != nil {
				res.Error(err.Error())
				return err
			}
		}
		err = tx.Delete(&bindProductList).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		err = tx.Where("contract_id = ?", contract.ID).Delete(&models.Hardware{}).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		err = tx.Delete(&contract).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		res.Success(DeleteSuccess, nil)
		return nil
	})
}

func GetContractList(c *gin.Context) {
	var res = NewResultMsg(c)
	var contracts []models.Contract
	err := dao.Db.Find(&contracts).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	for i, _ := range contracts {
		err = service.SyncCrm(&contracts[i])
		if err != nil {
			res.Error(err.Error())
			return
		}
		contracts[i].ContractProductList, err = service.GetContractProductList(contracts[i].ID)
		if err != nil {
			res.Error(err.Error())
			return
		}
	}
	res.Success(GetSuccess, List{List: contracts})
}

func GetContract(c *gin.Context) {
	var res = NewResultMsg(c)
	var contract models.Contract
	err := c.ShouldBind(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&contract).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = service.SyncCrm(&contract)
	if err != nil {
		res.Error(err.Error())
		return
	}
	contract.ContractProductList, err = service.GetContractProductList(contract.ID)
	if err != nil {
		res.Error(err.Error())
		return
	}
	contract.ContractHardwareList, err = service.GetContractHardwareList(contract.ID)
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, contract)
}

func GetCrmList(c *gin.Context) {
	var res = NewResultMsg(c)
	var crmList []models.ForPaasContractView
	err := dao.CrmDb.Find(&crmList).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, List{List: crmList})
}

//func PutContractJson(c *gin.Context){
//	var res = NewResultMsg(c)
//	var contract models.Contract
//	err:=c.ShouldBind(&contract)
//	if err != nil{
//		res.Error(err.Error())
//		return
//	}
//	err=dao.Db.First(&contract).Error
//	if err != nil{
//		res.Error(err.Error())
//		return
//	}
//	var binds []models.BindProduct
//	dao.Db.Where(&models.BindProduct{ContractID: contract.ID}).Find(&binds)
//	var products []models.Product
//	for _,bind:=range binds{
//		var product models.Product
//		dao.Db.First(&product,bind.Code)
//		err=json.Unmarshal([]byte(bind.FuncListJSON),&product.ContractFuncList)
//		if err!=nil{
//			res.Error(err.Error())
//		}
//		products=append(products,product)
//	}
//
//	var hosts []models.Hardware
//	dao.Db.Where(&models.Hardware{ContractID: contract.ID}).Find(&hosts)
//
//	productsJson,_:=json.Marshal(products)
//	hostsJson,_:=json.Marshal(hosts)
//
//	err=dao.Db.Model(&contract).Updates(models.Contract{ProductListJson: string(productsJson), HardwareListJson: string(hostsJson)}).Error
//	if err != nil{
//		res.Error(err.Error())
//		return
//	}
//	res.Success(PutSuccess,nil)
//}
