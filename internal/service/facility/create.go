package facility

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/net/context"
)

func (s FacilityService) CreateFacility(name string, parts []string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return 0, err
	}

	for _, part := range parts {
		err = validation.Validate(part, validation.Required)
		if err != nil {
			s.logger.Error("error: %v", err)
			return 0, fmt.Errorf("поле части помещения не должно быть пустым")
		}
	}

	id, err := s.facilityRepo.Create(ctx, name, parts)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, nil
}
