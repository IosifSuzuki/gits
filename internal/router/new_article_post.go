package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) NewArticlePOST(ctx *gin.Context) {
	log := r.container.GetLogger()

	accountValue, ok := ctx.Get(constant.AccountAppKey)
	if !ok {
		log.Error("can't retrieve account model from context")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, ok := accountValue.(*dto.Account)
	if !ok {
		log.Error("can't cast to account model from context value")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var newArticleForm dto.NewArticle
	if err := ctx.Bind(&newArticleForm); err != nil {
		log.Error("bind form has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := r.mainController.PostNewArticle(account, &newArticleForm); err != nil {
		log.Error("post new article has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/articles")
}
