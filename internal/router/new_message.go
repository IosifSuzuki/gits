package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"gits/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

func (r *router) NewMessage(ctx *gin.Context) {
	log := r.container.GetLogger()

	var newMessageForm dto.NewMessage
	if err := ctx.ShouldBind(&newMessageForm); err != nil {
		log.Error("bind form has failed", zap.Error(err))
		fieldErrs := handleMessageFormError(err)
		contactMe := dto.ContactMe{
			FieldErrors:   fieldErrs,
			FullName:      newMessageForm.FullName,
			Email:         newMessageForm.Email,
			Message:       newMessageForm.Message,
			FormSubmitted: false,
		}

		ctx.Redirect(http.StatusSeeOther, prepareRedirectLocationContactMe(&contactMe))
		return
	}

	if err := r.mainController.NewMessage(&newMessageForm); err != nil {
		log.Error("post message has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	contactMe := dto.ContactMe{
		FormSubmitted: true,
	}
	ctx.Redirect(http.StatusSeeOther, prepareRedirectLocationContactMe(&contactMe))
}

func prepareRedirectLocationContactMe(contactMe *dto.ContactMe) string {
	location := "/contact/me"

	base64ContactMeText, err := utils.MarshalQueryParam(contactMe)
	if err != nil {
		return location
	} else if base64ContactMeText == nil {
		return location
	}

	queryParams := url.Values{}
	queryParams.Set(constant.ReturnValueQueryParameter, *base64ContactMeText)
	urlLocation := url.URL{Path: location, RawQuery: queryParams.Encode()}

	return urlLocation.RequestURI()
}

func handleMessageFormError(err error) []dto.FieldError {
	var ve validator.ValidationErrors
	ok := errors.As(err, &ve)
	if !ok {
		return nil
	}

	fieldErrs := make([]dto.FieldError, 0, len(ve))
	for _, vi := range ve {
		fieldErr := getFieldError(vi)
		if fieldErr != nil {
			fieldErrs = append(fieldErrs, *fieldErr)
		}
	}

	return fieldErrs
}
