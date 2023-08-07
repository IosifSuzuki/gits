package form

type Authentication struct {
	Username string `form:"username" binding:"required,username"`
	Password string `form:"password" binding:"required,password"`
}
