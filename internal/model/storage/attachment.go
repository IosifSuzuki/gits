package storage

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	Path      *string
	Reference *string
}

func (a *Attachment) TableName() string {
	return "attachment"
}
