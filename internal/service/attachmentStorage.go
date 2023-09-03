package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gits/internal/container"
	"go.uber.org/zap"
	"io"
)

type AttachmentStorage interface {
	UploadAttachment(reader io.Reader, path string) (*string, error)
}

type attachmentStorage struct {
	sess      *session.Session
	container container.Container
}

func NewAttachmentStorage(container container.Container) (AttachmentStorage, error) {
	conf := container.GetConfig()
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(conf.AWS.S3Region),
		},
	})
	if err != nil {
		return nil, err
	}
	return &attachmentStorage{
		sess:      sess,
		container: container,
	}, nil
}

func (a *attachmentStorage) UploadAttachment(reader io.Reader, path string) (*string, error) {
	log := a.container.GetLogger()
	conf := a.container.GetConfig()
	uploader := s3manager.NewUploader(a.sess)
	uploadInput := s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(conf.AWS.S3Bucket),
		Key:    aws.String(path),
		ACL:    aws.String("public-read"),
	}
	uploadOutput, err := uploader.Upload(&uploadInput)
	if err != nil {
		log.Error("failed to upload image to s3 bucket")
		return nil, err
	}
	log.Debug("upload to s3 bucket image", zap.String("path", path))

	return &uploadOutput.Location, nil
}
