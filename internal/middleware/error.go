package middleware

import (
	"github.com/gin-gonic/gin"
	"gits/internal/container"
)

type ErrorHandler interface {
	Error() gin.HandlerFunc
}

type errorHandler struct {
	container.Container
}

func NewErrorHandler(container container.Container) ErrorHandler {
	return &errorHandler{
		Container: container,
	}
}

func (e *errorHandler) Error() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		httpCode := ctx.Writer.Status()
		if httpCode < 200 || httpCode >= 300 {
			ctx.HTML(httpCode, "views/error.tmpl", gin.H{
				"httpCode": httpCode,
			})
			return
		}
	}
}
