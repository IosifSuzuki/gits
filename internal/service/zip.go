package service

import (
	"archive/zip"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	"go.uber.org/zap"
	"io"
	"regexp"
	"strings"
)

type DecompressorFile interface {
	ExtractZip(request *app.ArticleUploadRequest) (*app.ArticleFiles, error)
}

type decompressorFile struct {
	container container.Container
}

func NewDecompressorFile(container container.Container) DecompressorFile {
	return &decompressorFile{
		container: container,
	}
}

func (d *decompressorFile) ExtractZip(request *app.ArticleUploadRequest) (*app.ArticleFiles, error) {
	log := d.container.GetLogger()

	zipReader, err := zip.NewReader(request.ReaderAt, request.Size)
	if err != nil {
		log.Error("create new zip reader has failed", zap.Error(err))
		return nil, err
	}
	var attachmentFiles = make(map[string]io.Reader)
	var mdFile io.Reader
	filterPathRegexp := regexp.MustCompile(`^[a-zA-z-_]+/([a-zA-z-_]+\.\w+)$`)
	for _, file := range zipReader.File {
		if !filterPathRegexp.MatchString(file.Name) {
			log.Debug("file not passed template path", zap.String("filename", file.Name))
			continue
		}
		matches := filterPathRegexp.FindStringSubmatch(file.Name)
		if len(matches) < 2 {
			err = errs.UnknownFileNameError
			log.Error("unknown file name", zap.Error(err))
			return nil, err
		}
		fileName := matches[1]
		f, err := file.Open()
		if err != nil {
			log.Error("cannot open file", zap.Error(err))
			return nil, err
		}
		if strings.HasSuffix(file.Name, ".md") {
			mdFile = f
		} else {
			attachmentFiles[fileName] = f
		}
	}
	return app.NewArticleFiles(mdFile, attachmentFiles), nil
}
