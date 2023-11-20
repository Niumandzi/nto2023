package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"golang.org/x/net/context"
	"regexp"
)

func (s EventService) CreateEvent(event model.Event) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&event,
		validation.Field(&event.Name, validation.Required),
		validation.Field(&event.Date, validation.Required, validation.By(validateDate)),
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

var basicDateRegex = regexp.MustCompile(`^(\d{2})\.(\d{2})\.(\d{4})$`)
