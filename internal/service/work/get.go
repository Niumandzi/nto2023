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
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return types, nil
}
