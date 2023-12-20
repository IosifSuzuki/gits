package router

import (
	"github.com/gin-gonic/gin"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	"gits/internal/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (r *router) ContactMe(ctx *gin.Context) {
	log := r.container.GetLogger()

	var contactMe dto.ContactMe
	returnValue := ctx.Query(constant.ReturnValueQueryParameter)
	if err := utils.Unmarshal(returnValue, &contactMe); err != nil {
		log.Error("fail map query param", zap.Error(err))
	}

	htmlContactMe, err := r.mainController.ContactMe(&contactMe)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "views/contact.tmpl", gin.H{
		"now":       time.Now(),
		"contactMe": htmlContactMe,
	})
}
