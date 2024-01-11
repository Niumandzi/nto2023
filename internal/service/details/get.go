package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

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

func (s DetailsService) GetActiveDetails(categoryName string) ([]model.Details, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.detailsRepo.GetActive(ctx, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.Details{}, err
	}

	return types, nil
}
