package dto

type NewMessage struct {
	FullName *string `json:"full_name" form:"full_name" binding:"required,full_name"`
	Email    *string `json:"email" form:"email" binding:"required,email"`
	Message  *string `json:"message" form:"message" binding:"required,filled_in=10"`
}
