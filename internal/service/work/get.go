package work

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s WorkTypeService) GetWorkTypes() ([]model.WorkType, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	types, err := s.workTypeRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (s WorkTypeService) GetActiveWorkTypes(categoryName string, facilityID int, status string) ([]model.WorkType, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	types, err := s.workTypeRepo.GetActive(ctx, categoryName, facilityID, status)
	if err != nil {
		return nil, err
	}

	return types, nil
}
