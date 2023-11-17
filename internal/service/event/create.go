package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"golang.org/x/net/context"
)

func (s EventService) CreateEvent(event model.EventWithDetails) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&event,
		validation.Field(&event.Name, validation.Required),
		validation.Field(&event.Details.TypeName, validation.Required),
		validation.Field(&event.Details.Category, validation.Required, validation.In("entertainment", "enlightenment", "education")),
	)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, err
	}

	detailId, err := s.detailsRepo.GetId(ctx, event.Details.Category, event.Details.TypeName)
	if err != nil {
		return 0, err
	}

	eventDB := model.Event{
		Id:          0,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		DetailsId:   detailId,
	}

	id, err := s.eventRepo.Create(ctx, eventDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s EventService) CreateDetails(categoryName string, typeName string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	id, err := s.detailsRepo.Create(ctx, categoryName, typeName)
	if err != nil {
		return 0, err
	}

	return id, err
}
