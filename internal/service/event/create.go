package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"golang.org/x/net/context"
)

func (s EventService) CreateEvent(event model.EventWithCategoryAndType) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&event,
		validation.Field(&event.Name, validation.Required),
		validation.Field(&event.Category.Category, validation.Required, validation.In("entertainment", "enlightenment", "education")),
		validation.Field(&event.Category.EventType.TypeName, validation.Required),
	)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, err
	}

	detailId, err := s.categoryTypeRepo.GetCategoryTypeID(ctx, event.Category.Category, event.Category.EventType.TypeName)

	event.Category.ID = detailId

	id, err := s.eventRepo.Create(ctx, event)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s EventService) CreateType(eventType string, eventCategory string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	id, err := s.categoryTypeRepo.CreateCategoryWithType(ctx, eventCategory, eventType)
	if err != nil {
		return 0, err
	}

	return id, err
}
