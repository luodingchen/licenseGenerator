package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"licenseGenerator/service"
	"os"
	"time"
)

func PostTemporaryLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}

	var contract models.Contract
	err = dao.Db.First(&contract, license.ContractID).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	if contract.ContractStatus == 0 {
		res.Error("Contract is closed")
		return
	}
	if contract.ContractTrailLicenseStatus == true {
		res.Error("This contract has already a trail license")
		return
	}
	if license.LicenseDuration > 30 {
		res.Error("Trail maximum duration is 30 days ")
		return
	}
	if license.LicenseDuration < 1 {
		res.Error("Trail minimum duration is 1 day")
		return
	}
	var key = service.GenerateKey()
	license.KeyID = key.ID
	license.LicensePublicKey = key.PublicKey
	license.LicenseType = 0
	license.LicenseStatus = 1

	contract.ContractTrailLicenseStatus = true

	dao.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Updates(&contract).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		err = tx.Create(&license).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		res.Success(PostSuccess, license)
		return nil
	})
}

func PostPermanentLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	var contract models.Contract
	dao.Db.First(&contract, license.ContractID)

	if contract.ContractStatus == 0 {
		res.Error("Contract is closed")
		return
	}
	if contract.ContractFormalLicenseStatus == true {
		res.Error("This contract has already a formal license")
		return
	}

	var key = service.GenerateKey()
	license.KeyID = key.ID
	license.LicensePublicKey = key.PublicKey
	license.LicenseType = 1

	contract.ContractFormalLicenseStatus = true

	dao.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Updates(&contract).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		err = tx.Create(&license).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		res.Success(PostSuccess, license)
		return nil
	})
}

func SignLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	if license.LicenseType == 1 {
		if license.LicenseStatus == 0 {
			res.Error("license has not got authorized")
			return
		}
		if license.LicenseStatus == 2 {
			res.Error("license has already got signed")
			return
		}
		if license.LicenseStatus == 3 {
			res.Error("license has been rejected by administrator")
			return
		}
		license.LicenseDuration = 36135
		t := time.Now()
		license.LicenseStartTime = t.Format("2006-01-02 15:04:05.000")
		license.LicenseDeadline = t.AddDate(99, 0, 0).Format("2006-01-02 15:04:05.000")
	}

	if license.LicenseType == 0 {
		t := time.Now()
		license.LicenseStartTime = time.Now().Format("2006-01-02 15:04:05.000")
		license.LicenseDeadline = t.AddDate(0, 0, license.LicenseDuration).Format("2006-01-02 15:04:05.000")
	}

	contractId := license.ContractID
	license.Config.ProductList, err = service.GetContractProductList(contractId)
	if err != nil {
		res.Error(err.Error())
		return
	}
	license.Config.HardwareList, err = service.GetContractHardwareList(contractId)
	if err != nil {
		res.Error(err.Error())
		return
	}

	license.Config.StartTime = license.LicenseStartTime
	license.Config.Deadline = license.LicenseDeadline

	configBytes, _ := json.Marshal(license.Config)

	var key models.RsaKey
	dao.Db.First(&key, license.KeyID)
	privateKey, err := service.DecodePrivateKeyString(key.PrivateKey)
	if err != nil {
		res.Error(err.Error())
		return
	}

	license.LicenseSignature, err = service.Sign(privateKey, configBytes)
	if err != nil {
		res.Error(err.Error())
		return
	}
	license.ConfigJson = string(configBytes)
	license.LicenseStatus = 2

	err = dao.Db.Updates(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PostSuccess, license)
}

func DeleteLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Delete(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(DeleteSuccess, nil)
}

func PutLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	if license.LicenseStatus == 3 {
		license.LicenseRejectedTime = time.Now().Format("2006-01-02 15:04:05.000")
	}
	err = dao.Db.Updates(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PutSuccess, license)
}

func DownloadLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	if license.LicenseStatus != 2 {
		res.Error("license has not got signed")
		return
	}
	fileBytes, _ := yaml.Marshal(license)

	var aesKey models.AesKey
	dao.Db.Last(&aesKey)
	encryptedBytes, err := service.AesEncrypt(fileBytes, []byte(aesKey.AesKeyString))

	err = ioutil.WriteFile("files/license", encryptedBytes, 0666)
	if err != nil {
		res.Error(err.Error())
		return
	}
	c.Header("Content-LicenseType", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+"license")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("files/license")
	err = os.Remove("files/license")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetLicense(c *gin.Context) {
	var res = NewResultMsg(c)
	var license models.License
	err := c.ShouldBind(&license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&license).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	if license.ConfigJson != "" {
		err = json.Unmarshal([]byte(license.ConfigJson), &license.Config)
		if err != nil {
			res.Error(err.Error())
			return
		}
	}
	res.Success(GetSuccess, license)
}

func GetLicenseList(c *gin.Context) {
	var res = NewResultMsg(c)
	var licenses []models.License
	err := dao.Db.Find(&licenses).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	for i, _ := range licenses {
		if licenses[i].ConfigJson != "" {
			err = json.Unmarshal([]byte(licenses[i].ConfigJson), &licenses[i].Config)
			if err != nil {
				res.Error(err.Error())
				return
			}
		}
		err = dao.Db.First(&licenses[i].Contract, licenses[i].ContractID).Error
		if err != nil {
			res.Error(err.Error())
			return
		}
	}
	res.Success(GetSuccess, List{List: licenses})
}
