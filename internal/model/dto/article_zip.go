package dto

import "io"

type ArticleUploadRequest struct {
	ReaderAt io.ReaderAt
	Size     int64
}
