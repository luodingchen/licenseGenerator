package controllers

import (
	"github.com/gin-gonic/gin"
	"licenseGenerator/dao"
	"licenseGenerator/models"
)

func PostFunc(c *gin.Context) {
	var res = NewResultMsg(c)
	var function models.Func
	err := c.ShouldBind(&function)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Create(&function).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PostSuccess, function)
}

func PostFuncList(c *gin.Context) {
	var res = NewResultMsg(c)
	var functionList []models.Func
	err := c.ShouldBind(&functionList)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Create(&functionList).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PostSuccess, List{List: functionList})
}

func PutFunc(c *gin.Context) {
	var res = NewResultMsg(c)
	var function models.Func
	err := c.ShouldBind(&function)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Updates(&function).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PutSuccess, function)
}

func DeleteFunc(c *gin.Context) {
	var res = NewResultMsg(c)
	var function models.Func
	err := c.ShouldBind(&function)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Delete(&function).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(DeleteSuccess, nil)
}
