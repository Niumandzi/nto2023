package facility

import (
	"context"
	"errors"
	"github.com/niumandzi/nto2023/model"
	"time"
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

func (s FacilityService) GetFacilitiesByDate(startDate string, startTime string, endDate string, endTime string) ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	start, _ := time.Parse("2006-01-02 15:04", startDate+" "+startTime)
	end, _ := time.Parse("2006-01-02 15:04", endDate+" "+endTime)

	if start.After(end) {
		err := errors.New("start date and time must be earlier than or equal to end date and time")
		s.logger.Error("Date and time range error: %v", err)
		return nil, err
	}

	facilities, err := s.facilityRepo.GetByDate(ctx, startDate, startTime, endDate, endTime)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return facilities, nil
}
