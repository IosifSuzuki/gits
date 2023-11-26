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
	service.Storage
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
	Articles() ([]*html.PreviewArticle, error)
	NewCategory(account *dto.Account) (*html.NewCategory, error)
	CreateNewCategory(account *dto.Account, form *dto.FormCategory) error
	ViewActions() ([]html.Action, error)
}

func NewMainController(
	container container.Container,
	storage service.Storage,
	session service.AccountSession,
	decompressorFile service.DecompressorFile,
	attachmentStorage service.AttachmentStorage,
	md service.MD,
	storageDAO storage.DAO,
) MainController {
	return &mainController{
		Container:         container,
		Storage:           storage,
		AccountSession:    session,
		DecompressorFile:  decompressorFile,
		AttachmentStorage: attachmentStorage,
		MD:                md,
		storageDAO:        storageDAO,
	}
}
