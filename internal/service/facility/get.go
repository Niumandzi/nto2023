package facility

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
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

func (s FacilityService) GetFacilitiesByDate(startDate string, endDate string) ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.Validate(startDate, validation.Required, validation.Date("2006-01-02"))
	if err != nil {
		s.logger.Error("Invalid start date format: %v", err)
		return nil, err
	}

	err = validation.Validate(endDate, validation.Required, validation.Date("2006-01-02"))
	if err != nil {
		s.logger.Error("Invalid end date format: %v", err)
		return nil, err
	}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	if start.After(end) {
		err = errors.New("start date must be earlier than or equal to end date")
		s.logger.Error("Date range error: %v", err)
		return nil, err
	}

	facilities, err := s.facilityRepo.GetByDate(ctx, startDate, endDate)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return facilities, nil
}
