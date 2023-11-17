package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) GetEventsByCategory(categoryName string) ([]model.EventWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, categoryName, "")
	if err != nil {
		return []model.EventWithDetails{}, err
	}

	return events, nil
}

func (s EventService) GetEventsByCategoryAndType(eventCategory string, eventType string) ([]model.EventWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, eventCategory, eventType)
	if err != nil {
		return []model.EventWithDetails{}, err
	}

	return events, nil
}

func (s EventService) GetDetailsByCategory(categoryName string) ([]model.Details, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.detailsRepo.Get(ctx, categoryName)
	if err != nil {
		return []model.Details{}, err
	}

	return types, nil
}
