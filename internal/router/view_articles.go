package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) ViewArticles(ctx *gin.Context) {
	log := r.container.GetLogger()

	var pageRequest dto.Page
	pageRequest.Page = retrievePage(ctx)

	articles, err := r.mainController.ViewArticles(&pageRequest)
	if err != nil {
		log.Error("view articles has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/articles.tmpl", articles)
}
