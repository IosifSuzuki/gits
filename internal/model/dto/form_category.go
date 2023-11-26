package dto

type FormCategory struct {
	Title *string `form:"title" binding:"required"`
}
