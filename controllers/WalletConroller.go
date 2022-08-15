package controllers

import (
	"ArvanWallet/resources"
	"ArvanWallet/services"
	"ArvanWallet/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUserBalance(c *gin.Context) {
	output := utils.NewOutput()

	mobileStr := c.Param("mobile")
	mobile, e := strconv.ParseInt(mobileStr, 10, 64)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	s := services.UserService{
		ServiceDB: services.DB,
		Redis:     services.REDIS,
		Ctx:       services.CTX,
	}
	balance, e := s.GetBalance(mobile)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusExpectationFailed, http.StatusExpectationFailed)
		return
	}

	out := resources.GetUserBalanceResource{
		Balance: *balance,
	}

	output.Data = out
	c.JSON(http.StatusOK, output)
	return
}

func GetUserTransactions(c *gin.Context) {
	output := utils.NewOutput()

	mobileStr := c.Param("mobile")
	mobile, e := strconv.ParseInt(mobileStr, 10, 64)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	s := services.UserService{
		ServiceDB: services.DB,
		Redis:     services.REDIS,
		Ctx:       services.CTX,
	}
	transactions, e := s.GetUserTransactions(mobile)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusExpectationFailed, http.StatusExpectationFailed)
		return
	}

	output.Data = transactions
	c.JSON(http.StatusOK, output)
	return
}
