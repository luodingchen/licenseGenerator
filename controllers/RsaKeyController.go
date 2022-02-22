package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"licenseGenerator/dao"
	"licenseGenerator/models"
)

func PostRsaKey(c *gin.Context) {
	var res = NewResultMsg(c)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	publicKey := &privateKey.PublicKey

	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	X509PublicKey := x509.MarshalPKCS1PublicKey(publicKey)

	privateKeyString := base64.StdEncoding.EncodeToString(X509PrivateKey)
	publicKeyString := base64.StdEncoding.EncodeToString(X509PublicKey)

	var key = models.RsaKey{
		PublicKey:  publicKeyString,
		PrivateKey: privateKeyString,
	}
	err := dao.Db.Create(&key).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PostSuccess, ID{ID: key.ID})
}

func GetRsaKey(c *gin.Context) {
	var res = NewResultMsg(c)
	var key []models.RsaKey
	err := c.ShouldBind(&key)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&key).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, key)
}

func GetRsaKeyList(c *gin.Context) {
	var res = NewResultMsg(c)
	var keys []models.RsaKey
	err := dao.Db.Find(&keys).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, List{List: keys})
}
