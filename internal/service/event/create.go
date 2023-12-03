package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"golang.org/x/net/context"
)

func (s EventService) CreateEvent(event model.Event) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&event,
		validation.Field(&event.Name, validation.Required),
		validation.Field(&event.Date, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&event.DetailsID, validation.Required, validation.Min(1).Error("Не выбран тип мероприятия")),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	eventDB := model.Event{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		DetailsID:   event.DetailsID,
	}

	id, err := s.eventRepo.Create(ctx, eventDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}
