package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (r *router) DashboardActions(ctx *gin.Context) {
	log := r.container.GetLogger()

	var pageRequest dto.Page
	pageRequest.Page = retrievePage(ctx)

	actions, err := r.mainController.ViewActions(&pageRequest)
	if err != nil {
		log.Error("view actions has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/actions.tmpl", actions)
}

func retrievePage(ctx *gin.Context) uint {
	page, ok := ctx.Params.Get("page")
	if !ok {
		return 1
	} else {
		number, err := strconv.Atoi(page)
		if err != nil {
			return 1
		} else if number <= 0 {
			return 1
		} else {
			return uint(number)
		}
	}
}
