package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"licenseGenerator/dao"
	"licenseGenerator/models"
)

func PostProduct(c *gin.Context) {
	var res = NewResultMsg(c)
	var product models.Product
	err := c.ShouldBind(&product)
	if err != nil {
		res.Error(err.Error())
		return
	}

	dao.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&product).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		if product.ProductFuncList != nil {
			for i, _ := range product.ProductFuncList {
				product.ProductFuncList[i].ProductID = product.ID
			}
			err = tx.Create(&product.ProductFuncList).Error
			if err != nil {
				res.Error(err.Error())
				return err
			}
		}
		res.Success(PostSuccess, product)
		return nil
	})
}

func PutProduct(c *gin.Context) {
	var res = NewResultMsg(c)
	var product models.Product
	err := c.ShouldBind(&product)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Updates(&product).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	if product.ProductFuncList != nil {
		for i, _ := range product.ProductFuncList {
			err = dao.Db.Updates(&product.ProductFuncList[i]).Error
			if err != nil {
				res.Error(err.Error())
				return
			}
		}
	}
	res.Success(PutSuccess, product)
}

func PutProductStatus(c *gin.Context) {
	var res = NewResultMsg(c)
	var product models.Product
	err := c.ShouldBind(&product)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Model(&product).Update("product_status", product.ProductStatus).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(PutSuccess, product)
}

func DeleteProduct(c *gin.Context) {
	var res = NewResultMsg(c)
	var product models.Product
	err := c.ShouldBind(&product)
	if err != nil {
		res.Error(err.Error())
		return
	}
	dao.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("product_id = ?", product.ID).Delete(&models.Func{}).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		err = tx.Delete(&product).Error
		if err != nil {
			res.Error(err.Error())
			return err
		}
		res.Success(DeleteSuccess, nil)
		return nil
	})
}

func GetProduct(c *gin.Context) {
	var res = NewResultMsg(c)
	var product models.Product
	err := c.ShouldBind(&product)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.First(&product).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = dao.Db.Where("product_id = ?", product.ID).Find(&product.ProductFuncList).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, product)
}

func GetProductList(c *gin.Context) {
	var res = NewResultMsg(c)
	var products []models.Product
	err := dao.Db.Find(&products).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Success(GetSuccess, List{List: products})
}

func GetUnclosedProductList(c *gin.Context) {
	var res = NewResultMsg(c)
	var products []models.Product
	err := dao.Db.Where("product_status = ?", 0).Find(&products).Error
	if err != nil {
		res.Error(err.Error())
		return
	}
	for i, product := range products {
		err = dao.Db.Where("product_id = ?", product.ID).Find(&products[i].ProductFuncList).Error
		if err != nil {
			res.Error(err.Error())
			return
		}
	}
	res.Success(GetSuccess, List{List: products})
}

//func DownloadProduct(c *gin.Context){
//	var res = NewResultMsg(c)
//	var product models.Product
//	err:=c.ShouldBind(&product)
//	if err != nil{
//		res.Error(err.Error())
//		return
//	}
//	err=dao.Db.First(&product).Error
//	if err != nil {
//		res.Error(err.Error())
//		return
//	}
//	fileDir:=product.Link
//	fileName:=product.Name
//	_, errByOpenFile := os.Open(fileDir + "/" + fileName)
//	if errByOpenFile != nil {
//		res.Error(err.Error())
//		return
//	}
//	c.Header("Content-LicenseType", "application/octet-stream")
//	c.Header("Content-Disposition", "attachment; filename="+fileName)
//	c.Header("Content-Transfer-Encoding", "binary")
//	c.File(fileDir + "/" + fileName)
//}
//
//func UploadProduct(c *gin.Context){
//	var res = NewResultMsg(c)
//	file,err:=c.FormFile("file")
//	if err!=nil{
//		res.Error(err.Error())
//		return
//	}
//	fileName:=filepath.Base(file.Filename)
//	err=c.SaveUploadedFile(file,"productFiles"+"/"+fileName)
//	if err!=nil{
//		res.Error(err.Error())
//		return
//	}
//	res.Success("Upload success",nil)
//}
