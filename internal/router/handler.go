package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gits/internal/container"
	"gits/internal/controller/central"
	"gits/internal/middleware"
	"gits/internal/model/dto"
	"gits/internal/utils"
	valid "gits/internal/validator"
	"html/template"
	"strings"
)

type Router interface {
	SetupHandlers() error
	Index(c *gin.Context)
}

type router struct {
	container              container.Container
	mainController         central.MainController
	authMiddleware         middleware.Auth
	errorHandlerMiddleware middleware.ErrorHandler
	observerMiddleware     middleware.Observer
}

func NewRouter(
	container container.Container,
	mainController central.MainController,
	authMiddleware middleware.Auth,
	errorHandlerMiddleware middleware.ErrorHandler,
	observerMiddleware middleware.Observer,
) Router {
	return &router{
		container:              container,
		mainController:         mainController,
		authMiddleware:         authMiddleware,
		errorHandlerMiddleware: errorHandlerMiddleware,
		observerMiddleware:     observerMiddleware,
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
	h.Use(r.observerMiddleware.Observer())
	h.GET("/", r.Index)
	h.GET("/articles/", r.Index)
	h.GET("/articles/:page", r.Index)
	h.GET("/article/:id", r.Article)
	h.POST("/contact/me", r.NewMessage)
	h.GET("/contact/me", r.ContactMe)
	h.GET("/auth", r.Authentication)
	h.POST("/auth", r.AuthenticationPOST)
	h.GET("/foxhole/artillery-calculator", r.FoxholeArtilleryCalculator)
	adminPanel := h.Group("/admin")
	adminPanel.Use(r.authMiddleware.Authorization())
	{
		adminPanel.GET("/new/article", r.NewArticle)
		adminPanel.POST("/new/article", r.NewArticlePOST)
		adminPanel.GET("/new/category", r.NewCategory)
		adminPanel.GET("/view/categories", r.ViewCategories)
		adminPanel.GET("/view/categories/:page", r.ViewCategories)
		adminPanel.GET("/view/articles/", r.ViewArticles)
		adminPanel.GET("/view/articles/:page", r.ViewArticles)
		adminPanel.POST("/new/category", r.NewCategoryPOST)
		adminPanel.GET("/actions/:page", r.DashboardActions)
		adminPanel.GET("/actions/", r.DashboardActions)
	}
	r.registerTemplateFunction(h)
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
		if err := engine.RegisterValidation("full_name", valid.FullNameValidator); err != nil {
			return err
		}
		if err := engine.RegisterValidation("filled_in", valid.FilledIn); err != nil {
			return err
		}
	}
	return nil
}

func (r *router) registerTemplateFunction(e *gin.Engine) {
	e.FuncMap = template.FuncMap{
		"DateFormat": utils.DateFormat,
		"Add":        utils.Add,
		"Mul":        utils.Mul,
	}
}

func getFieldError(err validator.FieldError) *dto.FieldError {
	var message string
	var field = strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		message = fmt.Sprintf("%s is required", field)
	case "min":
		message = fmt.Sprintf("%v doesn't contains enougth characters", field)
	default:
		message = fmt.Sprintf("%s has incorrect value", field)
	}

	return &dto.FieldError{
		Field:   utils.String(err.Field()),
		Message: utils.String(message),
	}
}
