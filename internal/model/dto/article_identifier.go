package dto

type ArticleIdentifier struct {
	ID int `uri:"id" binding:"required"`
}
