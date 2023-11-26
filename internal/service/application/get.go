package application

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationService) GetApplications(categoryName string, facilityId int, workTypeId int, status string) ([]model.ApplicationWithDetails, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	events, err := s.applicationRepo.Get(ctx, categoryName, facilityId, workTypeId, status)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return events, nil
}
