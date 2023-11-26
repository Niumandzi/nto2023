package facility

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityService) GetFacilities(categoryName string, workTypeID int, status string) ([]model.Facility, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	types, err := s.facilityRepo.Get(ctx, categoryName, workTypeID, status)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return types, nil
}
