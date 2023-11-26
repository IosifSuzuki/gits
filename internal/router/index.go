package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *router) Index(ctx *gin.Context) {
	log := r.container.GetLogger()

	htmlArticles, err := r.mainController.Articles()
	if err != nil {
		log.Error("fetch articles has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.HTML(http.StatusOK, "views/index.tmpl", gin.H{
		"now":      time.Now(),
		"articles": htmlArticles,
	})
}
