package services

import (
	"ArvanWallet/models"
	"ArvanWallet/requests"
	"context"
	"errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type UserService struct {
	ServiceDB *gorm.DB
	Redis     *redis.Client
	Ctx       context.Context
}

func (service UserService) GetBalance(userId int64) (*int64, error) {
	res, _ := service.Redis.Get(service.Ctx, strconv.FormatInt(userId, 10)).Result()
	if len(res) > 0 {
		balance, e := strconv.ParseInt(res, 10, 64)
		if e == nil {
			return &balance, nil
		}
	}

	var user models.Users
	if e := service.ServiceDB.Find(&user, userId).Error; e != nil {
		return nil, e
	}
	if user.Mobile > 0 {
		return &user.Balance, nil
	}
	return nil, errors.New("UserNotFound")
}

func (service UserService) GetUserTransactions(userId int64) (*[]models.UserTransactions, error) {
	var transactions []models.UserTransactions
	if e := service.ServiceDB.Order("id desc").
		Where("user_id = ?", userId).
		Find(&transactions).
		Error; e != nil {
		return nil, e
	}
	if transactions == nil {
		return nil, errors.New("TransactionsNotFound")
	}
	return &transactions, nil
}

func (service UserService) AddTransaction(input requests.AddTransactionRequest) error {
	service.Redis.Set(service.Ctx, strconv.FormatInt(input.Mobile, 10),
		input.Amount, time.Hour*2)

	trx := service.ServiceDB.Begin()

	var user models.Users
	if e := trx.Find(&user, input.Mobile).Error; e != nil {
		trx.Rollback()
		return e
	}
	if user.Mobile == 0 {
		user = models.Users{
			Mobile:    input.Mobile,
			Balance:   input.Amount,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if e := trx.Create(&user).Error; e != nil {
			trx.Rollback()
			return e
		}
	} else {
		user.Balance += input.Amount
		if e := trx.Save(&user).Error; e != nil {
			trx.Rollback()
			return e
		}
	}

	transaction := models.UserTransactions{
		UserId:     input.Mobile,
		Amount:     input.Amount,
		Reason:     input.Reason,
		NewBalance: user.Balance,
		CreatedAt:  time.Now(),
	}
	if e := trx.Create(&transaction).Error; e != nil {
		trx.Rollback()
		return e
	}

	trx.Commit()
	return nil
}
