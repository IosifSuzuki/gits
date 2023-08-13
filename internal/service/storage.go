package service

import (
	"encoding/base64"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"gits/internal/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage interface {
	AccountByCredential(credential *app.Credential) (*app.Account, error)
	AvailableCategories() ([]*app.Category, error)
	CreateNewArticle(article *app.NewArticle) (bool, error)
	RetrieveArticle(request *app.ArticleRequest) (*app.Article, error)
	RetrieveArticles() ([]*app.Article, error)
}

type storage struct {
	container.Container
	storageProvider provider.Storage
}

func NewStorage(container container.Container, storageProvider provider.Storage) Storage {
	return &storage{
		Container:       container,
		storageProvider: storageProvider,
	}
}

func (s *storage) AccountByCredential(credential *app.Credential) (*app.Account, error) {
	conn := s.storageProvider.GetConnection()

	var account stor.Account
	if err := conn.Where("username = ?", credential.Username).First(&account).Error; err != nil {
		return nil, err
	}
	decodedHashPassword, err := base64.StdEncoding.DecodeString(account.HashPassword)
	if err != nil {
		return nil, err
	}
	match, err := utils.CompareHashAndPassword(string(decodedHashPassword), credential.Password)
	if err != nil {
		return nil, err
	} else if !match {
		return nil, errs.NotMatchCredentialsError
	}
	var appAccount = app.Account{
		Id:           int(account.ID),
		Username:     account.Username,
		HashPassword: account.HashPassword,
		Role:         app.ParseString(string(account.Role)),
	}
	return &appAccount, nil
}

func (s *storage) AvailableCategories() ([]*app.Category, error) {
	conn := s.storageProvider.GetConnection()

	var categories []stor.Category

	if err := conn.Find(&categories).Error; err != nil {
		return nil, err
	}
	appCategories := make([]*app.Category, 0, len(categories))
	for _, category := range categories {
		appCategories = append(appCategories, &app.Category{
			Id:    int(category.ID),
			Title: category.Title,
		})
	}
	return appCategories, nil
}

func (s *storage) CreateNewArticle(article *app.NewArticle) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	storageCategories := make([]*stor.Category, 0, len(article.SelectedCategories))
	for _, categoryId := range article.SelectedCategories {
		storageCategories = append(storageCategories, &stor.Category{
			Model: gorm.Model{
				ID: uint(categoryId),
			},
			Title: nil,
		})
	}
	storageAttachments := make([]*stor.Attachment, 0, len(article.Attachments))
	for _, attachment := range article.Attachments {
		storageAttachments = append(storageAttachments, &stor.Attachment{
			Path: attachment.Path,
		})
	}
	storageArticle := stor.Article{
		PublisherId: article.PublisherId,
		Title:       article.Title,
		ReadingTime: article.ReadingTime,
		Location:    article.Location,
		Content:     article.Content,
		Categories:  storageCategories,
		Attachments: storageAttachments,
	}

	if err := conn.Save(&storageArticle).Error; err != nil {
		log.Error("can't save article to storage", zap.Error(err))
		return false, nil
	}
	return true, nil
}

func (s *storage) RetrieveArticle(request *app.ArticleRequest) (*app.Article, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	var article stor.Article
	if err := conn.Preload("Categories").First(&article, request.Id).Error; err != nil {
		log.Error("can't retrieve article by id", zap.Int("article id", request.Id))
		return nil, err
	}
	categories := make([]app.Category, 0)
	for _, category := range article.Categories {
		categories = append(categories, app.Category{
			Id:    int(category.ID),
			Title: category.Title,
		})
	}
	var account stor.Account
	if err := conn.First(&account, article.PublisherId).Error; err != nil {
		log.Error("can't retrieve account by id", zap.Error(err))
		return nil, err
	}
	author := app.Account{
		Id:           int(account.ID),
		Username:     account.Username,
		HashPassword: account.HashPassword,
		Role:         app.ParseString(string(account.Role)),
	}
	articleApp := app.Article{
		Id:          int(article.ID),
		Account:     &author,
		Title:       article.Title,
		Content:     article.Content,
		Location:    article.Location,
		ReadingTime: article.ReadingTime,
		Categories:  categories,
		UpdatedAt:   &article.UpdatedAt,
	}
	return &articleApp, nil
}

func (s *storage) RetrieveArticles() ([]*app.Article, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	storageArticles := make([]*stor.Article, 0)
	if err := conn.Find(&storageArticles).Error; err != nil {
		log.Error("find articles has failed", zap.Error(err))
		return nil, err
	}
	articles := make([]*app.Article, 0, len(storageArticles))
	for _, storageArticle := range storageArticles {
		categories := make([]app.Category, 0, len(storageArticle.Categories))
		for _, storageCategory := range storageArticle.Categories {
			categories = append(categories, app.Category{
				Id:    int(storageCategory.ID),
				Title: storageCategory.Title,
			})
		}
		articles = append(articles, &app.Article{
			Id:          int(storageArticle.ID),
			Account:     nil,
			Title:       storageArticle.Title,
			Content:     storageArticle.Content,
			Location:    storageArticle.Location,
			ReadingTime: storageArticle.ReadingTime,
			Categories:  categories,
			UpdatedAt:   &storageArticle.UpdatedAt,
		})
	}
	return articles, nil
}
