package central

import (
	"gits/internal/model/html"
	"go.uber.org/zap"
)

func (m *mainController) ViewActions() ([]html.Action, error) {
	log := m.GetLogger()

	storObservables, err := m.storageDAO.GetObservableRepository().RetrieveObservables()
	if err != nil {
		log.Error("retrieve observables has failed", zap.Error(err))
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

	return actions, err
}
