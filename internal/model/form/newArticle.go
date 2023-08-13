package form

import "mime/multipart"

type NewArticle struct {
	Title       *string               `form:"title" binding:"required"`
	Location    *string               `form:"location" binding:"required"`
	ReadingTime int                   `form:"reading_estimate" binding:"required"`
	Categories  []int                 `form:"categories" binding:"required"`
	ZipFile     *multipart.FileHeader `form:"zip_file" binding:"required"`
}
