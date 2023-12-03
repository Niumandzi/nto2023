package facility

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityService) GetFacilities() ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	facilities, err := s.facilityRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return facilities, nil
}

func (s FacilityService) GetActiveFacilities(categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	facilities, err := s.facilityRepo.GetActive(ctx, categoryName, workTypeID, status)
	if err != nil {
		return nil, err
	}

	return facilities, nil
}

func (s FacilityService) GetFacilitiesByDate(startDate string, endDate string) ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	facilities, err := s.facilityRepo.GetByDate(ctx, startDate, endDate)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return facilities, nil
}
