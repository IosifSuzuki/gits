package central

import (
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	stor "gits/internal/model/storage"
	"gits/internal/utils"
	"go.uber.org/zap"
	"strings"
)

func (m *mainController) NewMessage(message *dto.NewMessage) error {
	log := m.GetLogger()

	message.Message = utils.String(strings.TrimSpace(*message.Message))
	if len(*message.Message) == 0 {
		err := errs.NilError
		log.Error("message has empty string", zap.Error(err))
		return err
	}

	storMessage := stor.Message{
		FullName: message.FullName,
		Email:    message.Email,
		Message:  message.Message,
	}
	if _, err := m.storageDAO.GetMessageRepository().Create(&storMessage); err != nil {
		log.Error("create message has failed", zap.Error(err))
		return err
	}

	return nil
}
