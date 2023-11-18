package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

//func (s EventService) GetEventsByCategory(categoryName string) ([]model.EventWithDetails, error) {
//	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
//
//	defer cancel()
//
//	events, err := s.eventRepo.Get(ctx, categoryName, "")
//	if err != nil {
//		s.logger.Error("error: %v", err.Error())
//		return []model.EventWithDetails{}, err
//	}
//
//	return events, nil
//}

func (s DetailsService) GetDetails(categoryName string) ([]model.Details, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.detailsRepo.Get(ctx, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.Details{}, err
	}

	return types, nil
}
