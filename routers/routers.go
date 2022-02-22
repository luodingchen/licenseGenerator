package routers

import (
	"github.com/gin-gonic/gin"
	"licenseGenerator/controllers"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	//router.LoadHTMLGlob("template/*")
	//
	//router.GET("/", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.html", nil)
	//})
	//router.GET("/product.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "product.html", nil)
	//})
	//router.GET("/contract.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "contract.html", nil)
	//})
	//router.GET("/bind.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "bind.html", nil)
	//})
	//router.GET("/key.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "key.html", nil)
	//})
	//router.GET("/certification.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "certification.html", nil)
	//})
	//router.GET("/verify.html", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "verify.html", nil)
	//})

	productRouter := router.Group("product")
	{
		productRouter.POST("post", controllers.PostProduct)
		productRouter.GET("get", controllers.GetProduct)
		productRouter.GET("get_list", controllers.GetProductList)
		productRouter.GET("get_unclosed_list", controllers.GetUnclosedProductList)
		productRouter.PUT("put", controllers.PutProduct)
		productRouter.PUT("put_status", controllers.PutProductStatus)
		productRouter.DELETE("delete", controllers.DeleteProduct)
	}

	funcRouter := router.Group("func")
	{
		funcRouter.POST("post", controllers.PostFunc)
		funcRouter.POST("post_list", controllers.PostFuncList)
		funcRouter.PUT("put", controllers.PutFunc)
		funcRouter.DELETE("delete", controllers.DeleteFunc)
	}

	contractRouter := router.Group("contract")
	{
		contractRouter.POST("post", controllers.PostContract)
		contractRouter.GET("get", controllers.GetContract)
		contractRouter.GET("get_list", controllers.GetContractList)
		contractRouter.GET("get_crm_list", controllers.GetCrmList)
		contractRouter.PUT("put", controllers.PutContract)
		contractRouter.DELETE("delete", controllers.DeleteContract)
	}

	licenseRouter := router.Group("license")
	{
		licenseRouter.POST("post_temporary", controllers.PostTemporaryLicense)
		licenseRouter.POST("post_permanent", controllers.PostPermanentLicense)
		licenseRouter.GET("get", controllers.GetLicense)
		licenseRouter.GET("get_list", controllers.GetLicenseList)
		licenseRouter.GET("download", controllers.DownloadLicense)
		licenseRouter.PUT("put_status", controllers.PutLicense)
		licenseRouter.PUT("sign", controllers.SignLicense)
		licenseRouter.DELETE("delete", controllers.DeleteLicense)
	}

	aesKeyRouter := router.Group("aes_key")
	{
		aesKeyRouter.POST("post", controllers.PostAesKey)
	}

	rsaKeyRouter := router.Group("rsa_key")
	{
		rsaKeyRouter.POST("post", controllers.PostRsaKey)
		rsaKeyRouter.GET("get", controllers.GetRsaKey)
		rsaKeyRouter.GET("get_list", controllers.GetRsaKeyList)
	}

	router.POST("/verify_signature", controllers.VerifySignature)

	return router
}
