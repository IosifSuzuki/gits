package central

import (
	"gits/internal/model/html"
	"go.uber.org/zap"
)

// TODO refactor code

func (m *mainController) ViewActions() ([]html.Action, error) {
	log := m.GetLogger()

	observables, err := m.RetrieveObservables()
	if err != nil {
		log.Error("retrieve observables has failed", zap.Error(err))
		return nil, err
	}
	actions := make([]html.Action, 0, len(observables))
	for _, observable := range observables {
		action, err := html.NewAction(observable)
		if err != nil {
			log.Error("create new action has failed", zap.Error(err))
		} else if action != nil {
			actions = append(actions, *action)
		}
	}
	return actions, err
}
