package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *router) Article(ctx *gin.Context) {
	log := r.container.GetLogger()

	var articleRequest dto.ArticleIdentifier
	if err := ctx.ShouldBindUri(&articleRequest); err != nil {
		log.Error("bind uri has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	htmlArticle, err := r.mainController.Article(&articleRequest)
	if err != nil {
		log.Error("article by id has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/article.tmpl", gin.H{
		"now":     time.Now(),
		"article": htmlArticle,
	})
}