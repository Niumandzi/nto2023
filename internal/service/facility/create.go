package facility

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/net/context"
)

func (s FacilityService) CreateFacility(name string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return 0, err
	}

	id, err := s.facilityRepo.Create(ctx, name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, err
}
