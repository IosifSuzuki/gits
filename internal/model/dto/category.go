package dto

import (
	stor "gits/internal/model/storage"
)

type Category struct {
	Id    int
	Title *string
}

func NewCategory(storCategory stor.Category) *Category {
	return &Category{
		Id:    int(storCategory.ID),
		Title: storCategory.Title,
	}
}
