package storage

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	PublisherId int
	Title       *string
}

func (c *Category) TableName() string {
	return "category"
}
