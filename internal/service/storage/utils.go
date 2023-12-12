package storage

import "gorm.io/gorm"

func Pagination(page int, batch int) func(db *gorm.DB) *gorm.DB {
	offset := (page - 1) * batch

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(batch)
	}
}
