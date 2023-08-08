package html

import "gits/internal/model/response"

type NewArticle struct {
	PublisherName       string
	AvailableCategories []response.Category
}
