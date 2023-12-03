package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) GetEvents(categoryName string, detailsID int) ([]model.EventWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, categoryName, detailsID)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.EventWithDetails{}, err
	}

	return events, nil
}

func (s EventService) GetActiveEvents(categoryName string) ([]model.EventWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.GetActive(ctx, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.EventWithDetails{}, err
	}

	return events, nil
}
