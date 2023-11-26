package central

import (
	strip "github.com/grokify/html-strip-tags-go"
	"gits/internal/model/html"
	"gits/internal/utils"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) Articles() ([]*html.PreviewArticle, error) {
	log := m.GetLogger()

	storArticles, err := m.storageDAO.GetArticleRepository().RetrieveArticles()
	if err != nil {
		log.Error("retrieve articles from storage has failed", zap.Error(err))
		return nil, err
	}

	previewArticles := make([]*html.PreviewArticle, 0, len(storArticles))
	for _, storArticle := range storArticles {
		stripHtmlContent := strip.StripTags(*storArticle.Content)
		content := utils.PrefixString(stripHtmlContent, 70)
		contentHTML := template.HTML(content)
		date := storArticle.UpdatedAt
		previewArticles = append(previewArticles, &html.PreviewArticle{
			Id:      int(storArticle.ID),
			Title:   storArticle.Title,
			Date:    &date,
			Content: &contentHTML,
		})
	}
	return previewArticles, nil
}
