package central

import (
	"encoding/base64"
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"gits/internal/model/html"
	stor "gits/internal/model/storage"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) Article(articleIdentifier *dto.ArticleIdentifier) (*html.Article, error) {
	log := m.GetLogger()

	storArticle, err := m.storageDAO.GetArticleRepository().Article(articleIdentifier.ID)
	if err != nil {
		log.Error("fail to retrieve article by id", zap.Error(err))

		return nil, err
	}

	htmlContent, err := m.convertMdToHtmlContent(storArticle)
	if err != nil {
		log.Error("convert mark down to html has failed", zap.Error(err))

		return nil, err
	}

	storAccount, err := m.storageDAO.GetAccountRepository().AccountByIdentifier(storArticle.PublisherId)
	if err != nil {
		log.Error("fail to retrieve publisher from article", zap.Error(err))

		return nil, err
	} else if storAccount == nil {
		err = errs.NilError
		log.Error("publisher contains nil value", zap.Error(err))

		return nil, err
	}

	categories := make([]dto.Category, 0, len(storArticle.Categories))
	for _, storCategory := range storArticle.Categories {
		category := dto.NewCategory(storCategory)
		categories = append(categories, *category)
	}

	account := dto.NewAccount(*storAccount)
	htmlArticle := html.Article{
		Author:      account,
		Title:       storArticle.Title,
		Location:    storArticle.Location,
		Date:        &storArticle.UpdatedAt,
		ReadingTime: storArticle.ReadingTime,
		Content:     template.HTML(*htmlContent),
		Categories:  categories,
	}
	return &htmlArticle, nil
}

func (m *mainController) convertMdToHtmlContent(storArticle *stor.Article) (*string, error) {
	log := m.GetLogger()

	mdData, err := m.articleData(storArticle)
	if err != nil {
		log.Error("fail to get article data", zap.Error(err))
		return nil, err
	}

	htmlData, err := m.MD.RenderMdToHTML(mdData, m.getAttachmentIdentifiers(storArticle))
	if err != nil {
		log.Error("rendering md to html has failed", zap.Error(err))

		return nil, err
	}

	content := string(htmlData)

	return &content, nil
}

func (m *mainController) convertMdToHtmlPreview(storArticle *stor.Article, words uint) (*string, error) {
	log := m.GetLogger()

	mdData, err := m.articleData(storArticle)
	if err != nil {
		log.Error("fail to get article data", zap.Error(err))
		return nil, err
	}

	htmlData, err := m.MD.RenderMdToPreviewHTML(mdData, words)
	if err != nil {
		log.Error("rendering md to html has failed", zap.Error(err))

		return nil, err
	}

	content := string(htmlData)

	return &content, nil
}

func (m *mainController) getAttachmentIdentifiers(article *stor.Article) map[string]string {
	log := m.GetLogger()
	identifiers := make(map[string]string)

	for _, attachment := range article.Attachments {
		ref := attachment.Reference
		path := attachment.Path
		if ref == nil {
			log.Error("article's reference contains nil")
			continue
		}
		if path == nil {
			log.Error("article's path contains nil")
			continue
		}

		identifiers[*attachment.Reference] = *attachment.Path
	}

	return identifiers
}

func (m *mainController) articleData(storArticle *stor.Article) ([]byte, error) {
	log := m.GetLogger()

	mdContent := storArticle.Content
	if mdContent == nil {
		err := errs.NilError
		log.Error("article absent content", zap.Error(err))

		return nil, err
	}

	mdData, err := base64.StdEncoding.DecodeString(*mdContent)
	if err != nil {
		log.Error("base64 decoding md text failed", zap.Error(err))

		return nil, err
	}

	return mdData, err
}
