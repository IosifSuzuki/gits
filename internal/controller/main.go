package controller

import (
	"context"
	"github.com/google/uuid"
	strip "github.com/grokify/html-strip-tags-go"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	"gits/internal/model/form"
	"gits/internal/model/html"
	"gits/internal/model/request"
	"gits/internal/model/response"
	"gits/internal/service"
	"gits/internal/utils"
	"go.uber.org/zap"
	"html/template"
	"io"
	"strings"
)

type mainController struct {
	container.Container
	service.Storage
	service.AccountSession
	service.DecompressorFile
	service.AttachmentStorage
	service.MD
}

type MainController interface {
	Auth(authentication *form.Authentication) (*response.AuthSessionResponse, error)
	NewArticle(account *app.Account) (*html.NewArticle, error)
	PostNewArticle(account *app.Account, form *form.NewArticle) error
	Article(request *request.Article) (*html.Article, error)
	Articles() ([]*html.PreviewArticle, error)
	NewCategory(account *app.Account) (*html.NewCategory, error)
	CreateNewCategory(account *app.Account, form *form.NewCategory) error
}

func NewMainController(
	container container.Container,
	storage service.Storage,
	session service.AccountSession,
	decompressorFile service.DecompressorFile,
	attachmentStorage service.AttachmentStorage,
	md service.MD,
) MainController {
	return &mainController{
		Container:         container,
		Storage:           storage,
		AccountSession:    session,
		DecompressorFile:  decompressorFile,
		AttachmentStorage: attachmentStorage,
		MD:                md,
	}
}

func (m *mainController) Auth(authentication *form.Authentication) (*response.AuthSessionResponse, error) {
	log := m.GetLogger()
	ctx := context.Background()

	credential := &app.Credential{
		Username: authentication.Username,
		Password: authentication.Password,
	}
	account, err := m.Storage.AccountByCredential(credential)
	if err != nil {
		log.Error("retrieve account by credential has failed", zap.Error(err))
		return nil, err
	}
	accountSession, err := m.AccountSession.CreateAccountSession(ctx, account)
	if err != nil {
		log.Error("create session has failed", zap.Error(err))
		return nil, err
	}

	return &response.AuthSessionResponse{
		SessionId: accountSession.SessionId,
	}, nil
}

func (m *mainController) NewArticle(account *app.Account) (*html.NewArticle, error) {
	log := m.GetLogger()

	categories, err := m.Storage.AvailableCategories()
	if err != nil {
		log.Error("retrieve available categories has failed", zap.Error(err))
		return nil, err
	}
	resCategories := make([]response.Category, 0, len(categories))
	for _, category := range categories {
		resCategories = append(resCategories, response.Category{
			Id:    category.Id,
			Title: category.Title,
		})
	}
	return &html.NewArticle{
		PublisherName:       account.Username,
		AvailableCategories: resCategories,
	}, err
}

func (m *mainController) PostNewArticle(account *app.Account, form *form.NewArticle) error {
	log := m.GetLogger()

	file, err := form.ZipFile.Open()
	if err != nil {
		log.Error("open zip file has failed", zap.Error(err))
		return err
	}
	articleRequest := app.ArticleUploadRequest{
		ReaderAt: file,
		Size:     form.ZipFile.Size,
	}
	articleFiles, err := m.DecompressorFile.ExtractZip(&articleRequest)
	if err != nil {
		log.Error("cannot extract zip file", zap.Error(err))
		return err
	}
	for filename, _ := range articleFiles.Attachments {
		log.Debug("zip file contains", zap.String("filename", filename))
	}
	newAttachmentIdentifiers := make(map[string]string)
	folderName := strings.ToLower(*form.Title)
	folderName = strings.Join(strings.Split(folderName, " "), "_")
	for attachmentPath, attachmentReader := range articleFiles.Attachments {
		identifier, err := uuid.NewRandom()
		filePath := folderName + "/" + identifier.String()
		if err != nil {
			log.Error("create new uuid random has failed", zap.Error(err))
			return err
		}
		newPath, err := m.AttachmentStorage.UploadAttachment(attachmentReader, filePath)
		if err != nil {
			log.Error("upload attachment has failed", zap.Error(err))
			return err
		} else if newPath == nil {
			log.Error("upload attachment return empty path", zap.Error(err))
			return errs.NilError
		}
		newAttachmentIdentifiers[attachmentPath] = *newPath
	}
	mdFileData, err := io.ReadAll(articleFiles.MDFile)
	if err != nil {
		log.Error("read mark down file has failed", zap.Error(err))
		return err
	}
	transformedMDData, err := m.MD.RenderMdToHTML(mdFileData, newAttachmentIdentifiers)
	if err != nil {
		log.Error("transform md to new view has failed", zap.Error(err))
		return err
	}
	transformedMString := string(transformedMDData)
	attachments := make([]app.NewAttachment, 0, len(articleFiles.Attachments))
	for _, attachmentIdentifier := range newAttachmentIdentifiers {
		var newAttachmentIdentifier = attachmentIdentifier
		attachments = append(attachments, app.NewAttachment{
			Path: &newAttachmentIdentifier,
		})
	}
	newArticle := app.NewArticle{
		PublisherId:        account.Id,
		Title:              form.Title,
		Location:           form.Location,
		ReadingTime:        form.ReadingTime,
		SelectedCategories: form.Categories,
		Content:            &transformedMString,
		Attachments:        attachments,
	}
	ok, err := m.Storage.CreateNewArticle(&newArticle)
	if err != nil {
		log.Error("create new article has failed", zap.Error(err))
		return err
	} else if !ok {
		err = errs.UnsuccessfulCreateArticleError
		log.Error("create new article has failed", zap.Error(err))
		return err
	}
	return err
}

func (m *mainController) Article(request *request.Article) (*html.Article, error) {
	log := m.GetLogger()

	appArticleRequest := app.ArticleRequest{
		Id: request.Id,
	}
	article, err := m.RetrieveArticle(&appArticleRequest)
	if err != nil {
		log.Error("retrieve article by id has failed", zap.Error(err))
		return nil, err
	}
	htmlCategories := make([]html.Category, 0, len(article.Categories))
	for _, category := range article.Categories {
		htmlCategories = append(htmlCategories, html.Category{
			Id:    category.Id,
			Title: category.Title,
		})
	}
	author := html.Author{
		FullName: &article.Account.Username,
	}
	htmlArticle := html.Article{
		Author:      &author,
		Title:       article.Title,
		Location:    article.Location,
		Date:        article.UpdatedAt,
		ReadingTime: article.ReadingTime,
		Content:     template.HTML(*article.Content),
		Categories:  htmlCategories,
	}
	return &htmlArticle, nil
}

func (m *mainController) Articles() ([]*html.PreviewArticle, error) {
	log := m.GetLogger()

	articles, err := m.RetrieveArticles()
	if err != nil {
		log.Error("retrieve articles from storage has failed", zap.Error(err))
		return nil, err
	}
	previewArticles := make([]*html.PreviewArticle, 0, len(articles))
	for _, article := range articles {
		stripHtmlContent := strip.StripTags(*article.Content)
		content := utils.PrefixString(stripHtmlContent, 70)
		contentHTML := template.HTML(content)
		previewArticles = append(previewArticles, &html.PreviewArticle{
			Id:      article.Id,
			Title:   article.Title,
			Date:    article.UpdatedAt,
			Content: &contentHTML,
		})
	}
	return previewArticles, nil
}

func (m *mainController) NewCategory(account *app.Account) (*html.NewCategory, error) {
	log := m.GetLogger()

	if account == nil {
		err := errs.NilError
		log.Error("account has nil value", zap.Error(err))
		return nil, err
	}
	return &html.NewCategory{
		PublisherName: account.Username,
	}, nil
}

func (m *mainController) CreateNewCategory(account *app.Account, form *form.NewCategory) error {
	log := m.GetLogger()

	var newCategory = &app.NewCategory{
		PublisherId: account.Id,
		Title:       form.Title,
	}
	ok, err := m.Storage.CreateNewCategory(newCategory)
	if !ok {
		log.Error("create new category has completed unsuccessful")
		return nil
	} else if err != nil {
		log.Error("create new category has failed", zap.Error(err))
		return err
	}
	return nil
}
