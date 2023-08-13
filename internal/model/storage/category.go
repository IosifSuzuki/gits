package storage

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title *string
}

func (c *Category) TableName() string {
	return "category"
}
