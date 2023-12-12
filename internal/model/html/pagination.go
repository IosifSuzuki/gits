package html

type Pagination struct {
	PreviousItem *PaginationItem
	Items        []PaginationItem
	NextItem     *PaginationItem
}
