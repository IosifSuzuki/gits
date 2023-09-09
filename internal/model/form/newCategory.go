package form

type NewCategory struct {
	Title *string `form:"title" binding:"required"`
}
