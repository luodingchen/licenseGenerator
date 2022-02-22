package dao

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
)

var Db *gorm.DB
var CrmDb *gorm.DB

func InitAllMysql() (err error) {
	Db, err = InitMysql("configs/mysql.yaml")
	if err != nil {
		return
	}
	CrmDb, err = InitMysql("configs/crm.yaml")
	if err != nil {
		return
	}
	return nil
}

func InitMysql(filename string) (Db *gorm.DB, err error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("读取mysql文件失败")
		return
	}
	var da Database
	err = yaml.Unmarshal(fileBytes, &da)
	if err != nil {
		fmt.Println("解析mysql文件失败")
		return
	}
	//dsn := "root:root@tcp(116.62.132.128:3306)/licenseGenerator?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=%v&parseTime=True&loc=Local", da.Name, da.Password, da.Protocol, da.IP, da.Port, da.Library, da.CharSet)
	fmt.Println(dsn)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
