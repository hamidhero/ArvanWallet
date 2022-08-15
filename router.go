package main

import (
	"ArvanWallet/controllers"
	"github.com/gin-gonic/gin"
)

func GetRouter(port string) (router *gin.Engine) {
	router = gin.Default()
	router.ForwardedByClientIP = true
	router.RedirectFixedPath = true

	api := router.Group("api")
	{
		wallet := api.Group("wallet")
		{
			wallet.GET("balance/:mobile", controllers.GetUserBalance)
			wallet.GET("transactions/:mobile", controllers.GetUserTransactions)
		}
	}

	router.Run("0.0.0.0:" + port)
	return
}
