package central

import (
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/utils"
)

func (m *mainController) ContactMe(contactMe *dto.ContactMe) (*html.ContactMe, error) {
	return &html.ContactMe{
		FullName:             contactMe.FullName,
		Email:                contactMe.Email,
		Message:              contactMe.Message,
		FieldErrs:            contactMe.FieldErrors,
		FormSubmitted:        contactMe.FormSubmitted,
		FormSubmittedMessage: utils.String("form submitted successfully"),
	}, nil
}
