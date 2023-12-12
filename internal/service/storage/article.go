package storage

import (
	"gits/internal/container"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"go.uber.org/zap"
)

type ArticleRepository interface {
	CreateNewArticle(article *stor.Article) error
	Article(id int) (*stor.Article, error)
	Articles(page int, batch int) ([]stor.Article, error)
	LenArticles() (uint, error)
	AvailableCategories() ([]stor.Category, error)
	CreateNewCategory(category *stor.Category) error
	Categories(ids []int) ([]stor.Category, error)
}

type articleRepository struct {
	container       container.Container
	storageProvider provider.Storage
}

func NewArticleRepository(container container.Container, storageProvider provider.Storage) ArticleRepository {
	return &articleRepository{
		container:       container,
		storageProvider: storageProvider,
	}
}

func (a *articleRepository) AvailableCategories() ([]stor.Category, error) {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	var categories []stor.Category

	if err := conn.Find(&categories).Error; err != nil {
		log.Error("failed to retrieve all categories", zap.Error(err))
		return nil, err
	}

	return categories, nil
}

func (a *articleRepository) CreateNewArticle(article *stor.Article) error {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	if err := conn.Save(&article).Error; err != nil {
		log.Error("can't save article to storage", zap.Error(err))
		return err
	}
	return nil
}

func (a *articleRepository) Article(id int) (*stor.Article, error) {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	var article stor.Article
	if err := conn.Preload("Categories").Preload("Attachments").First(&article, id).Error; err != nil {
		log.Error("can't retrieve article by id", zap.Int("article id", id))
		return nil, err
	}

	return &article, nil
}

func (a *articleRepository) CreateNewCategory(category *stor.Category) error {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	if err := conn.Save(&category).Error; err != nil {
		log.Error("save category has failed", zap.Error(err))
		return err
	}
	return nil
}

func (a *articleRepository) Categories(ids []int) ([]stor.Category, error) {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	categories := make([]stor.Category, 0)
	if err := conn.Where(ids).Find(&categories).Error; err != nil {
		log.Error("unable to get categories by ids", zap.Error(err))
	}

	return categories, nil
}

func (a *articleRepository) LenArticles() (uint, error) {
	conn := a.storageProvider.GetConnection()
	log := a.container.GetLogger()

	var rows uint
	if err := conn.Model(stor.Article{}).Select("count(*)").Find(&rows).Error; err != nil {
		log.Error("retrieve articles has failed", zap.Error(err))
		return 0, err
	}

	return rows, nil
}

func (a *articleRepository) Articles(page int, batch int) ([]stor.Article, error) {
	log := a.container.GetLogger()
	conn := a.storageProvider.GetConnection()

	var articles []stor.Article
	query := conn
	if err := query.Scopes(Pagination(page, batch)).Find(&articles).Error; err != nil {
		log.Error("observables with pagination has failed", zap.Error(err))
		return nil, err
	}

	return articles, nil
}
