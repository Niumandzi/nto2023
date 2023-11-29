package facility

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityService) GetFacilities(categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	facilities, err := s.facilityRepo.Get(ctx, categoryName, workTypeID, status)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return nil, err
	}

	return facilities, nil
}
