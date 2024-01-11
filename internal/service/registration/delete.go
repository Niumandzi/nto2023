package registration

import "context"

func (r RegistrationService) DeleteRegistration(registrationID int) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.contextTimeout)
	defer cancel()

	err := r.registrationRepo.Delete(ctx, registrationID)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		return err
	}
	return nil
}
