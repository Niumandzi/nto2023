package registration

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (r RegistrationService) GetRegistrations(facilityID int, mugID int, teacherID int) ([]model.RegistrationWithDetails, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.contextTimeout)
	defer cancel()

	registrations, err := r.registrationRepo.Get(ctx, facilityID, mugID, teacherID)
	if err != nil {
		return nil, err
	}

	return registrations, nil
}
