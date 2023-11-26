package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *router) AuthenticationPOST(ctx *gin.Context) {
	log := r.container.GetLogger()
	conf := r.container.GetConfig()

	cookieTTL := int(conf.Cache.SessionTTL / time.Second)
	log.Debug("ttl cookie for session in second", zap.Int("ttl in second", cookieTTL))

	var authenticationForm dto.Authentication
	if err := ctx.Bind(&authenticationForm); err != nil {
		log.Error("bind form has failed", zap.Error(err))
		return
	}

	authSessionResponse, err := r.mainController.Auth(&authenticationForm)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.SetCookie(constant.CookieSessionKey, authSessionResponse.SessionId, cookieTTL, "/", "", false, false)

	ctx.Redirect(http.StatusSeeOther, "admin/new/article")
}
