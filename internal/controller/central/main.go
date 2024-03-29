package central

import (
	"gits/internal/container"
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/service"
	"gits/internal/service/storage"
)

type mainController struct {
	container.Container
	service.AccountSession
	service.DecompressorFile
	service.AttachmentStorage
	service.MD
	storageDAO storage.DAO
}

type MainController interface {
	Auth(authentication *dto.Authentication) (*dto.AccountSession, error)
	NewArticle(account *dto.Account) (*html.NewArticle, error)
	PostNewArticle(account *dto.Account, form *dto.NewArticle) error
	Article(articleIdentifier *dto.ArticleIdentifier) (*html.Article, error)
	Articles(page *dto.Page) (*html.Articles, error)
	NewCategory(account *dto.Account) (*html.NewCategory, error)
	CreateNewCategory(account *dto.Account, form *dto.FormCategory) error
	ViewActions(page *dto.Page) (*html.Actions, error)
	ContactMe(contactMe *dto.ContactMe) (*html.ContactMe, error)
	NewMessage(message *dto.NewMessage) error
	ViewCategories(page *dto.Page) (*html.Categories, error)
	ViewArticles(page *dto.Page) (*html.Articles, error)
}

func NewMainController(
	container container.Container,
	session service.AccountSession,
	decompressorFile service.DecompressorFile,
	attachmentStorage service.AttachmentStorage,
	md service.MD,
	storageDAO storage.DAO,
) MainController {
	return &mainController{
		Container:         container,
		AccountSession:    session,
		DecompressorFile:  decompressorFile,
		AttachmentStorage: attachmentStorage,
		MD:                md,
		storageDAO:        storageDAO,
	}
}
