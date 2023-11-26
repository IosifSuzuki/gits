package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) NewCategory(ctx *gin.Context) {
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

	newCategory, err := r.mainController.NewCategory(account)
	if err != nil {
		log.Error("new category operation has failed", zap.Error(err))
	}

	ctx.HTML(http.StatusOK, "views/newCategory.tmpl", newCategory)
}

