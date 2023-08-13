package html

import (
	"fmt"
	"html/template"
	"time"
)

type Article struct {
	Author      *Author
	Title       *string
	Location    *string
	Date        *time.Time
	ReadingTime int
	Content     template.HTML
	Categories  []Category
}

func (a *Article) ReadingTimeText() string {
	if a.ReadingTime < 60 {
		return fmt.Sprintf("%d m", a.ReadingTime)
	} else {
		return fmt.Sprintf("%d h %d m", a.ReadingTime/60, a.ReadingTime%60)
	}
}
