package storage

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	PublisherId int
	Title       *string
	ReadingTime int
	Location    *string
	Content     *string
	Categories  []Category   `gorm:"many2many:article_category"`
	Attachments []Attachment `gorm:"many2many:article_attachment"`
}

func (c *Article) TableName() string {
	return "article"
}
