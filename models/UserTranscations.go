package models

import "time"

type UserTransactions struct {
	Id         int64      `gorm:"column:id;primary_key" json:"id"`
	UserId     int64      `gorm:"column:user_id" json:"mobile"`
	Amount     int64      `gorm:"column:amount" json:"amount"`
	Reason     string     `gorm:"column:reason" json:"reason"`
	NewBalance int64      `gorm:"column:new_balance" json:"new_balance"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserTransactions) TableName() string {
	return "UserTransactions"
}
