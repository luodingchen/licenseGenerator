package main

import (
	"licenseGenerator/dao"
	"licenseGenerator/models"
	"licenseGenerator/routers"
)

func main() {
	err := dao.InitAllMysql()
	if err != nil {
		panic(err)
	}
	dao.Db.AutoMigrate(&models.AesKey{})
	dao.Db.AutoMigrate(&models.RsaKey{})
	dao.Db.AutoMigrate(&models.Product{})
	dao.Db.AutoMigrate(&models.Func{})
	dao.Db.AutoMigrate(&models.Contract{})
	dao.Db.AutoMigrate(&models.BindProduct{})
	dao.Db.AutoMigrate(&models.BindFunc{})
	dao.Db.AutoMigrate(&models.Hardware{})
	dao.Db.AutoMigrate(&models.License{})
	r := routers.InitRouters()
	r.Run(":80")
}
