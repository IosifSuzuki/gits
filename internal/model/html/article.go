package html

import (
	"fmt"
	"gits/internal/model/dto"
	"html/template"
	"time"
)

type Article struct {
	Author      *dto.Account
	Title       *string
	Location    *string
	Date        *time.Time
	ReadingTime int
	Content     template.HTML
	Categories  []dto.Category
}

func (a *Article) ReadingTimeText() string {
	if a.ReadingTime < 60 {
		return fmt.Sprintf("%d m", a.ReadingTime)
	} else {
		return fmt.Sprintf("%d h %d m", a.ReadingTime/60, a.ReadingTime%60)
	}
}
