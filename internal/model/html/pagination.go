package html

type Pagination struct {
	Page         int
	Batch        int
	PreviousItem *PaginationItem
	Items        []PaginationItem
	NextItem     *PaginationItem
}
