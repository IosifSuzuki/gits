package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (r *router) Authentication(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/auth.tmpl", gin.H{
		"now": time.Now(),
	})
}
