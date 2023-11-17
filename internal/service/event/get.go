package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) GetEventsByCategory(eventCategory string) ([]model.EventWithCategoryAndType, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, eventCategory, "")
	if err != nil {
		return []model.EventWithCategoryAndType{}, err
	}

	return events, nil
}

func (s EventService) GetEventsByCategoryAndType(eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, eventCategory, eventType)
	if err != nil {
		return []model.EventWithCategoryAndType{}, err
	}

	return events, nil
}

func (s EventService) GetTypesByCategory(eventCategory string) ([]model.EventType, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.categoryTypeRepo.GetCategoryTypes(ctx, eventCategory)
	if err != nil {
		return []model.EventType{}, err
	}

	return types, nil
}
