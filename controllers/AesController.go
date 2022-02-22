package controllers

import (
	"github.com/gin-gonic/gin"
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"licenseGenerator/service"
)

func PostAesKey(c *gin.Context) {
	var res = NewResultMsg(c)
	var aesKey models.AesKey
	aesKey.AesKeyString = service.RandString(16)
	dao.Db.Create(&aesKey)
	res.Success(PostSuccess, aesKey)
}

func GetLastAesKey(c *gin.Context) {
	var res = NewResultMsg(c)
	var aesKey models.AesKey
	dao.Db.Last(&aesKey)
	res.Success(GetSuccess, aesKey)
}
