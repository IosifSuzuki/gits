package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (r *router) FoxholeArtilleryCalculator(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/foxholeArtillery.tmpl", gin.H{
		"now": time.Now(),
	})
}
