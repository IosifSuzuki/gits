package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) NewArticle(ctx *gin.Context) {
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

	newArticle, err := r.mainController.NewArticle(account)
	if err != nil {
		log.Error("NewArticle has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Debug("account has gain access to create new post", zap.Int("accountId", account.ID))

	ctx.HTML(http.StatusOK, "views/newArticle.tmpl", newArticle)
}
