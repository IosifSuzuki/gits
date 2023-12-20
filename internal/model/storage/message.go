package storage

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	FullName *string
	Email    *string
	Message  *string
}

func (m *Message) TableName() string {
	return "message"
}
