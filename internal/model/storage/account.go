package storage

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Role         Role
	Username     string
	HashPassword string
}

func (a *Account) TableName() string {
	return "account"
}
