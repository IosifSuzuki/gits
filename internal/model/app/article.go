package app

import "time"

type Article struct {
	Id          int
	Account     *Account
	Title       *string
	Content     *string
	Location    *string
	ReadingTime int
	Categories  []Category
	UpdatedAt   *time.Time
}
