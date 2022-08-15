package tests

import (
	"ArvanWallet/requests"
	"ArvanWallet/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserBalance(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()

	s := services.UserService{
		ServiceDB: services.DB,
		Redis:     services.REDIS,
		Ctx:       services.CTX,
	}
	b, e := s.GetBalance(9376019072)
	assert.Equal(t, nil, e)
	assert.NotEqual(t, nil, b)
}

func TestGetUserTransactions(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()

	s := services.UserService{
		ServiceDB: services.DB,
		Redis:     services.REDIS,
		Ctx:       services.CTX,
	}
	b, e := s.GetUserTransactions(9376019072)
	assert.Equal(t, nil, e)
	assert.NotEqual(t, nil, b)
}

func TestAddTransaction(t *testing.T)  {
	services.ConnectDB()
	services.ConnectRedis()

	s := services.UserService{
		ServiceDB: services.DB,
		Redis:     services.REDIS,
		Ctx:       services.CTX,
	}

	req := requests.AddTransactionRequest{
		Mobile: 9376019072,
		Amount: 1000,
		Reason: "test",
	}

	e := s.AddTransaction(req)
	assert.Equal(t, nil, e)

	req.Amount = -1000
	e = s.AddTransaction(req)
	assert.Equal(t, nil, e)
}

func TestGetUserBalanceApi(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()
	services.ConnectRabbitMq()

	path := "/api/wallet/balance/9376019072"
	router := gin.Default()
	router.GET(path)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserTransactionsApi(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()
	services.ConnectRabbitMq()

	path := "/api/wallet/transactions/9376019072"
	router := gin.Default()
	router.GET(path)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}