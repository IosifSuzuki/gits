package central

import (
	"encoding/base64"
	"github.com/google/uuid"
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"gits/internal/model/storage"
	"go.uber.org/zap"
	"io"
	"strings"
)

func (m *mainController) PostNewArticle(account *dto.Account, form *dto.NewArticle) error {
	log := m.GetLogger()

	file, err := form.ZipFile.Open()
	if err != nil {
		log.Error("open zip file has failed", zap.Error(err))
		return err
	}

	articleZip := dto.ArticleUploadRequest{
		ReaderAt: file,
		Size:     form.ZipFile.Size,
	}
	articleFiles, err := m.DecompressorFile.ExtractZip(&articleZip)
	if err != nil {
		log.Error("cannot extract zip file", zap.Error(err))
		return err
	}

	for filename := range articleFiles.Attachments {
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

	base64MDContent := base64.StdEncoding.EncodeToString(mdFileData)

	storAttachments := make([]storage.Attachment, 0, len(articleFiles.Attachments))
	for reference, attachmentIdentifier := range newAttachmentIdentifiers {
		var pathAttachment = attachmentIdentifier
		var reference = reference
		storAttachment := storage.Attachment{
			Path:      &pathAttachment,
			Reference: &reference,
		}

		storAttachments = append(storAttachments, storAttachment)
	}

	selectedCategories, err := m.storageDAO.GetArticleRepository().Categories(form.Categories)
	if err != nil {
		log.Error("can't get categories by ids", zap.Error(err))
		return err
	}

	storArticle := storage.Article{
		PublisherId: account.ID,
		Title:       form.Title,
		ReadingTime: form.ReadingTime,
		Location:    form.Location,
		Content:     &base64MDContent,
		Categories:  selectedCategories,
		Attachments: storAttachments,
	}

	err = m.storageDAO.GetArticleRepository().CreateNewArticle(&storArticle)
	if err != nil {
		log.Error("create new article has failed", zap.Error(err))
		return err
	}

	return err
}
