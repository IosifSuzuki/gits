package service

import (
	"gits/internal/container"
	"io"
)

type AttachmentStorage interface {
	UploadAttachment(reader io.Reader, path string) (bool, error)
}

type attachmentStorage struct {
}

func NewAttachmentStorage(container container.Container) AttachmentStorage {
	return &attachmentStorage{}
}

func (a *attachmentStorage) UploadAttachment(reader io.Reader, path string) (bool, error) {
	return true, nil
}
