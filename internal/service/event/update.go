package event

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/internal/validations"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) UpdateEvent(eventUpd model.Event) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&eventUpd,
		validation.Field(&eventUpd.Name, validation.Required),
		validation.Field(&eventUpd.Date, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&eventUpd.DetailsID, validation.Required, validation.Min(1).Error("Не выбран тип мероприятия")),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	err = s.eventRepo.Update(ctx, eventUpd)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
