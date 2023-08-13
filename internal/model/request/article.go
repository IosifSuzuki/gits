package request

type Article struct {
	Id int `uri:"id" binding:"required"`
}
