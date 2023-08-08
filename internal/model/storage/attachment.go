package storage

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	Path string
}
