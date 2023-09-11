package app

import "time"

type Observable struct {
	AccountId *int
	Account   *Account
	Ip        *Ip
	Browser   *string
	OS        *string
	OSVersion *string
	Path      *string
	Device    *string
	UpdatedAt *time.Time
}
