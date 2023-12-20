package html

import (
	"gits/internal/model/dto"
)

type ContactMe struct {
	FullName             *string
	Email                *string
	Message              *string
	FieldErrs            []dto.FieldError
	FormSubmitted        bool
	FormSubmittedMessage *string
}
