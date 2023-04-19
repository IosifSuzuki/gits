package controller

import (
	"github.com/gin-gonic/gin"
	"gits/internal/container"
	"gits/internal/services"
	"html/template"
	"net"
	"net/http"
	"time"
)

type Controller interface {
	SetupHandlers() error
	Index(c *gin.Context)
}

type controller struct {
	container container.Container
	publisher services.Publisher
}

func NewController(container container.Container, publisher services.Publisher) Controller {
	return &controller{
		container: container,
		publisher: publisher,
	}
}

func (c *controller) SetupHandlers() error {
	conf := c.container.GetConfig()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(c.Error)
	r.GET("/", c.Index)
	r.GET("/articles", c.Index)
	r.GET("/article", c.Article)
	r.LoadHTMLGlob("web/**/*")
	r.Static("/assets", "./assets")
	return r.Run(net.JoinHostPort(conf.Host, conf.Port))
}

func (c *controller) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "views/index.tmpl", gin.H{
		"now": time.Now(),
	})
}

func (c *controller) Article(ctx *gin.Context) {
	a, _ := c.publisher.Article()
	ctx.HTML(http.StatusOK, "views/article.tmpl", gin.H{
		"now":     time.Now(),
		"article": template.HTML(a.Content),
	})
}

func (c *controller) Error(ctx *gin.Context) {
	ctx.Next()
	httpCode := ctx.Writer.Status()
	if httpCode < 200 || httpCode >= 300 {
		ctx.HTML(httpCode, "views/error.tmpl", gin.H{
			"httpCode": httpCode,
		})
	}
}
