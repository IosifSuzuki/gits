package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) ViewCategories(ctx *gin.Context) {
	log := r.container.GetLogger()

	var pageRequest dto.Page
	pageRequest.Page = retrievePage(ctx)

	categories, err := r.mainController.ViewCategories(&pageRequest)
	if err != nil {
		log.Error("view categories has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/categories.tmpl", categories)
}
