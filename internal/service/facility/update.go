package facility

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s FacilityService) UpdateFacility(facilityId int, name string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.facilityRepo.Update(ctx, facilityId, name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
