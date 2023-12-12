package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *router) Index(ctx *gin.Context) {
	log := r.container.GetLogger()

	var pageRequest dto.Page
	pageRequest.Page = retrievePage(ctx)

	htmlArticles, err := r.mainController.Articles(&pageRequest)
	if err != nil {
		log.Error("fetch articles has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.HTML(http.StatusOK, "views/index.tmpl", gin.H{
		"now":     time.Now(),
		"content": htmlArticles,
	})
}
