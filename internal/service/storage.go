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
	CreateNewCategory(category *app.NewCategory) (bool, error)
	ExistsIp(ip string) (bool, error)
	CreateNewIp(ip *app.Ip) (*uint, error)
	RetrieveIp(ip string) (*app.Ip, error)
	CreateNewObservable(observable *app.Observable) (bool, error)
	RetrieveObservables() ([]app.Observable, error)
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

func (s *storage) CreateNewCategory(category *app.NewCategory) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	storageCategory := stor.Category{
		PublisherId: category.PublisherId,
		Title:       category.Title,
	}

	if err := conn.Save(&storageCategory).Error; err != nil {
		log.Error("save category has failed", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (s *storage) ExistsIp(ip string) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	var exists bool
	if err := conn.Model(stor.Ip{}).Select("count(*) > 0").Where("ip = ?", ip).Find(&exists).Error; err != nil {
		log.Error("checking is exist ip has failed", zap.Error(err))
		return false, err
	}
	return exists, nil
}

func (s *storage) RetrieveIp(ip string) (*app.Ip, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	var storageIp stor.Ip
	if err := conn.Where("ip = ?", ip).First(&storageIp).Error; err != nil {
		log.Error("retrieve ip model from db has failed", zap.Error(err))
		return nil, err
	}
	return &app.Ip{
		ID:       &storageIp.ID,
		Ip:       storageIp.Ip,
		Hostname: storageIp.Hostname,
		City:     storageIp.City,
		Region:   storageIp.Region,
		Country:  storageIp.Country,
		Loc:      storageIp.Loc,
		Org:      storageIp.Org,
		Postal:   storageIp.Postal,
		Timezone: storageIp.Timezone,
	}, nil
}

func (s *storage) CreateNewIp(ip *app.Ip) (*uint, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	storageIp := stor.Ip{
		Ip:         ip.Ip,
		Hostname:   ip.Hostname,
		City:       ip.City,
		Region:     ip.Region,
		Country:    ip.Country,
		Loc:        ip.Loc,
		Org:        ip.Org,
		Postal:     ip.Postal,
		Timezone:   ip.Timezone,
		Observable: nil,
	}

	if err := conn.Save(&storageIp).Error; err != nil {
		log.Error("save ip to db has failed", zap.Error(err))
		return nil, err
	}
	return &storageIp.ID, nil
}

func (s *storage) CreateNewObservable(observable *app.Observable) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	ipAppModel := observable.Ip
	if ipAppModel == nil {
		err := errs.NilError
		log.Error("ip has nil value", zap.Error(err))
		return false, err
	}
	existIpStorageModel, err := s.ExistsIp(*ipAppModel.Ip)
	if err != nil {
		log.Error("exists ip operation has failed", zap.Error(err))
		return false, err
	}
	if !existIpStorageModel {
		ipID, err := s.CreateNewIp(ipAppModel)
		if err != nil {
			log.Error("create new ip has failed", zap.Error(err))
			return false, err
		}
		ipAppModel.ID = ipID
	} else {
		ipAppModel, err = s.RetrieveIp(*ipAppModel.Ip)
		if err != nil {
			log.Error("retrieve ip has failed", zap.Error(err))
			return false, err
		}
	}
	storageObservable := stor.Observable{
		AccountId: observable.AccountId,
		IpId:      *ipAppModel.ID,
		Browser:   observable.Browser,
		OS:        observable.OS,
		OSVersion: observable.OSVersion,
		Path:      observable.Path,
		Device:    observable.Device,
	}

	if err := conn.Save(&storageObservable).Error; err != nil {
		log.Error("observable model has failed to save in db", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (s *storage) RetrieveObservables() ([]app.Observable, error) {
	log := s.GetLogger()
	conn := s.storageProvider.GetConnection()

	var observables []stor.Observable
	if err := conn.Order("updated_at desc").Preload("Ip").Preload("Account").Find(&observables).Error; err != nil {
		log.Error("retrieve observables has failed", zap.Error(err))
		return nil, err
	}
	var appObservables = make([]app.Observable, 0, len(observables))
	for _, observable := range observables {
		var (
			appAccount *app.Account
			appIp      *app.Ip
		)
		storageAccount := observable.Account
		if storageAccount != nil {
			appAccount = &app.Account{
				Id:           int(storageAccount.ID),
				Username:     storageAccount.Username,
				HashPassword: storageAccount.HashPassword,
				Role:         app.Role(storageAccount.Role),
			}
		}
		storageIp := observable.Ip
		if storageIp != nil {
			appIp = &app.Ip{
				ID:       &storageIp.ID,
				Ip:       storageIp.Ip,
				Hostname: storageIp.Hostname,
				City:     storageIp.City,
				Region:   storageIp.Region,
				Country:  storageIp.Country,
				Loc:      storageIp.Loc,
				Org:      storageIp.Org,
				Postal:   storageIp.Postal,
				Timezone: storageIp.Timezone,
			}
		}
		appObservables = append(appObservables, app.Observable{
			AccountId: observable.AccountId,
			Account:   appAccount,
			Ip:        appIp,
			Browser:   observable.Browser,
			OS:        observable.OS,
			OSVersion: observable.OSVersion,
			Path:      observable.Path,
			Device:    observable.Device,
			UpdatedAt: &observable.UpdatedAt,
		})
	}
	return appObservables, nil
}
