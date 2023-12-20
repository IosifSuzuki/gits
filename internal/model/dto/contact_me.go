package dto

type ContactMe struct {
	FullName      *string      `form:"full_name"`
	Email         *string      `form:"email"`
	Message       *string      `form:"message"`
	FieldErrors   []FieldError `json:"field_errors"`
	FormSubmitted bool         `json:"form_submitted"`
}
