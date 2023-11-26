package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) NewCategoryPOST(ctx *gin.Context) {
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
	}

	var newCategory dto.FormCategory
	if err := ctx.Bind(&newCategory); err != nil {
		log.Error("bind form has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := r.mainController.CreateNewCategory(account, &newCategory); err != nil {
		log.Error("create new category has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/admin/new/category")
}