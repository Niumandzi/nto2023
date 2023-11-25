package work

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/internal/validations"
	"github.com/niumandzi/nto2023/model"
)

func (w WorkService) CreateApplication(application model.Application) (int, error) {
	ctx, cancel := context.WithTimeout(w.ctx, w.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&application,
		validation.Field(&application.Description, validation.Required),
		validation.Field(&application.CreatedAt, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&application.DetailsID, validation.Required, validation.Min(1).Error("Не выбран тип мероприятия")),
	)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return 0, err
	}

	eventDB := model.Event{
		Name:        application.Name,
		Date:        application.Date,
		Description: application.Description,
		DetailsID:   application.DetailsID,
	}

	id, err := w.eventRepo.Create(ctx, eventDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}
