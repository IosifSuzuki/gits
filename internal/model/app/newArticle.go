package app

type NewArticle struct {
	PublisherId        int
	Title              *string
	Location           *string
	ReadingTime        int
	SelectedCategories []int
	Content            *string
	Attachments        []NewAttachment
}
