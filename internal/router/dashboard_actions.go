package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (r *router) DashboardActions(ctx *gin.Context) {
	log := r.container.GetLogger()

	actions, err := r.mainController.ViewActions()
	if err != nil {
		log.Error("view actions has failed", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/actions.tmpl", actions)
}
