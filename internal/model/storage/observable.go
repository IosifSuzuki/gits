package storage

import "gorm.io/gorm"

type Observable struct {
	gorm.Model
	AccountId *int
	Account   *Account
	IpId      uint
	Ip        *Ip
	Browser   *string
	OS        *string
	OSVersion *string
	Path      *string
	Device    *string
}

func (o *Observable) TableName() string {
	return "observable"
}
