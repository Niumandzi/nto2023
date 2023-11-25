package facility

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (s FacilityService) UpdateFacility(facility model.Facility) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.Validate(facility.Name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.facilityRepo.Update(ctx, facility.ID, facility.Name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
