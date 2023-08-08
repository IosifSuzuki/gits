package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gits/internal/container"
	"gits/internal/controller"
	"gits/internal/middleware"
	"gits/internal/model/app"
	"gits/internal/model/constant"
	"gits/internal/model/form"
	"gits/internal/service"
	valid "gits/internal/validator"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"time"
)

type Router interface {
	SetupHandlers() error
	Index(c *gin.Context)
}

type router struct {
	container              container.Container
	publisher              service.Publisher
	mainController         controller.MainController
	authMiddleware         middleware.Auth
	errorHandlerMiddleware middleware.ErrorHandler
}

func NewRouter(
	container container.Container,
	publisher service.Publisher,
	mainController controller.MainController,
	authMiddleware middleware.Auth,
	errorHandlerMiddleware middleware.ErrorHandler,
) Router {
	return &router{
		container:              container,
		publisher:              publisher,
		mainController:         mainController,
		authMiddleware:         authMiddleware,
		errorHandlerMiddleware: errorHandlerMiddleware,
	}
}

func (r *router) SetupHandlers() error {
	if err := r.registerValidators(); err != nil {
		return err
	}
	conf := r.container.GetConfig()
	h := gin.New()
	h.Use(gin.Logger())
	h.Use(gin.Recovery())
	h.Use(r.errorHandlerMiddleware.Error())
	h.GET("/", r.Index)
	h.GET("/articles", r.Index)
	h.GET("/article", r.Article)
	h.GET("/auth", r.Authentication)
	h.POST("/auth", r.AuthenticationPOST)
	adminPanel := h.Group("/admin")
	adminPanel.Use(r.authMiddleware.Authorization())
	{
		adminPanel.GET("/new/article", r.NewArticle)
	}
	h.LoadHTMLGlob("web/**/*")
	h.Static("/assets", "./assets")
	return h.Run(conf.Server.Address())
}

func (r *router) registerValidators() error {
	if engine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := engine.RegisterValidation("password", valid.PasswordValidator); err != nil {
			return err
		}
		if err := engine.RegisterValidation("username", valid.UsernameValidator); err != nil {
			return err
		}
	}
	return nil
}

func (r *router) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/index.tmpl", gin.H{
		"now": time.Now(),
	})
}

func (r *router) Article(ctx *gin.Context) {
	a, _ := r.publisher.Article()
	ctx.HTML(http.StatusOK, "views/article.tmpl", gin.H{
		"now":     time.Now(),
		"article": template.HTML(a.Content),
	})
}

func (r *router) Authentication(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/auth.tmpl", gin.H{
		"now": time.Now(),
	})
}

func (r *router) AuthenticationPOST(ctx *gin.Context) {
	log := r.container.GetLogger()
	conf := r.container.GetConfig()

	cookieTTL := int(conf.Cache.SessionTTL / time.Second)
	log.Debug("ttl cookie for session in second", zap.Int("ttl in second", cookieTTL))

	var authenticationForm form.Authentication
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

func (r *router) NewArticle(ctx *gin.Context) {
	log := r.container.GetLogger()

	accountValue, ok := ctx.Get(constant.AccountAppKey)
	if !ok {
		log.Error("can't retrieve account model from context")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	account, ok := accountValue.(*app.Account)
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
	log.Debug("account has gain access to create new post", zap.Int("accountId", account.Id))
	ctx.HTML(http.StatusOK, "views/newArticle.tmpl", newArticle)
}
