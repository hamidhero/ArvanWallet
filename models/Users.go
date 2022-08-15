package models

import "time"

type Users struct {
	Mobile    int64     `gorm:"column:mobile;primary_key"`
	Balance   int64     `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Users) TableName() string {
	return "Users"
}