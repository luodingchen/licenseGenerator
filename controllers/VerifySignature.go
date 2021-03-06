package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"licenseGenerator/service"
	"os"
	"path/filepath"
	"time"
)

func VerifySignature(c *gin.Context) {
	var res = NewResultMsg(c)
	file, err := c.FormFile("file")
	if err != nil {
		res.Error(err.Error())
		return
	}
	fileDir := "files"
	fileName := filepath.Base(file.Filename)
	err = c.SaveUploadedFile(file, fileDir+"/"+fileName)
	if err != nil {
		res.Error(err.Error())
		return
	}
	fs, e := ioutil.ReadFile(fileDir + "/" + fileName)
	if e != nil {
		res.Error(err.Error())
		return
	}
	err = os.Remove(fileDir + "/" + fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var aesKey models.AesKey
	dao.Db.Last(&aesKey)
	decrypted, _ := service.AesDecrypt(string(fs), aesKey.AesKeyString)
	var license models.License
	err = yaml.Unmarshal(decrypted, &license)
	if err != nil {
		res.Error(err.Error())
		return
	}
	configBytes := []byte(license.ConfigJson)

	err = json.Unmarshal([]byte(license.ConfigJson), &license.Config)
	if err != nil {
		res.Error(err.Error())
		return
	}

	publicKey, err := service.DecodePublicKeyString(license.LicensePublicKey)
	if err != nil {
		res.Error(err.Error())
		return
	}

	//privateKey, err := service.DecodePrivateKeyString(key.PrivateKey)
	//if err != nil {
	//	res.Error(err.Error())
	//	return
	//}

	err = service.Verify(publicKey, configBytes, license.LicenseSignature)
	if err != nil {
		res.Error(err.Error())
		return
	}
	deadlineTime, err := time.Parse("2006-01-02 15:04:05.000", license.LicenseDeadline)
	if deadlineTime.Before(time.Now()) {
		res.Error("license had expired")
		return
	}
	//
	//var hardwareInfoList []models.HardwareInfo
	//for _, hardware := range license.Config.HardwareList {
	//	var hardwareInfo models.HardwareInfo
	//	cpuBytes, err := service.Decrypt(publicKey, privateKey, hardware.Cpu)
	//	if err != nil {
	//		res.Error(err.Error())
	//		return
	//	}
	//	diskBytes, err := service.Decrypt(publicKey, privateKey, hardware.Disk)
	//	if err != nil {
	//		res.Error(err.Error())
	//		return
	//	}
	//	hostBytes, err := service.Decrypt(publicKey, privateKey, hardware.Host)
	//	if err != nil {
	//		res.Error(err.Error())
	//		return
	//	}
	//	netBytes, err := service.Decrypt(publicKey, privateKey, hardware.Net)
	//	if err != nil {
	//		res.Error(err.Error())
	//		return
	//	}
	//	json.Unmarshal(cpuBytes, &hardwareInfo.Cpu)
	//	json.Unmarshal(diskBytes, &hardwareInfo.Disk)
	//	json.Unmarshal(hostBytes, &hardwareInfo.Host)
	//	json.Unmarshal(netBytes, &hardwareInfo.Net)
	//	hardwareInfoList = append(hardwareInfoList, hardwareInfo)
	//}

	res.Success("LicenseSignature verification succeeded", license.Config)
}
