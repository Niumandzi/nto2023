package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) GetEvents(categoryName string, detailsID int, isActive bool) ([]model.EventWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.eventRepo.Get(ctx, categoryName, detailsID, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.EventWithDetails{}, err
	}

	return events, nil
}
