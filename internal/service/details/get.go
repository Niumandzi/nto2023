package event

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s DetailsService) GetDetails(categoryName string, isActive bool) ([]model.Details, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.detailsRepo.Get(ctx, categoryName, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.Details{}, err
	}

	return types, nil
}
