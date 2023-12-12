package central

import (
	"fmt"
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/service"
	"go.uber.org/zap"
	"html/template"
)

const BatchSize uint = 20

func (m *mainController) ViewActions(page *dto.Page) (*html.Actions, error) {
	log := m.GetLogger()

	pagination, err := m.prepareActionsPagination(page)
	if err != nil {
		log.Error("fail to prepare pagination for action", zap.Error(err))
		return nil, err
	}

	storObservables, err := m.storageDAO.GetObservableRepository().Observables(int(page.Page), int(BatchSize))
	if err != nil {
		log.Error("obtain observables by pagination has failed", zap.Error(err))
		return nil, err
	}

	actions := make([]html.Action, 0, len(storObservables))
	for _, storObservable := range storObservables {
		action, err := html.NewAction(&storObservable)
		if err != nil {
			log.Error("create new action has failed", zap.Error(err))
		} else if action != nil {
			actions = append(actions, *action)
		}
	}

	return &html.Actions{
		Actions:    actions,
		Pagination: pagination,
	}, err
}

func (m *mainController) prepareActionsPagination(page *dto.Page) (*html.Pagination, error) {
	log := m.GetLogger()

	countObservables, err := m.storageDAO.GetObservableRepository().LenObservables()
	if err != nil {
		log.Error("obtain count of observable has failed", zap.Error(err))
		return nil, err
	}

	paginationBuilder := service.NewPagination(page.Page, BatchSize, countObservables)

	paginationItemBuilderFunc := func(idx uint) *html.PaginationItem {
		path := fmt.Sprintf("/admin/actions/%d", idx)
		title := template.HTML(fmt.Sprintf("%d", idx))
		isActive := idx == page.Page

		return html.NewPaginationItem(path, title, isActive)
	}
	paginationItemNextBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("/admin/actions/%d", page.Page+1)
		title := template.HTML("&raquo")
		return html.NewPaginationItem(path, title, false)
	}
	paginationItemPreviousBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("/admin/actions/%d", page.Page-1)
		title := template.HTML("&laquo")
		return html.NewPaginationItem(path, title, false)
	}

	pagination := paginationBuilder.Build(
		paginationItemBuilderFunc,
		paginationItemNextBuilderFunc,
		paginationItemPreviousBuilderFunc,
	)

	return pagination, nil
}
