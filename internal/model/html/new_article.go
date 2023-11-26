package html

import (
	"gits/internal/model/dto"
)

type NewArticle struct {
	PublisherName       string
	AvailableCategories []dto.Category
}
