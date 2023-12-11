package dto

import "gits/internal/model/storage"

type Account struct {
	ID           int
	Username     string
	HashPassword string
	Role         Role
}

func NewAccount(storAccount storage.Account) *Account {
	return &Account{
		ID:           int(storAccount.ID),
		Username:     storAccount.Username,
		HashPassword: storAccount.HashPassword,
		Role:         ParseString(storAccount.Role.String()),
	}
}

func (a *Account) FullName() string {
	return a.Username
}
