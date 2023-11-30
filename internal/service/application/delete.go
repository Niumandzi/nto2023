package application

import "context"

func (s ApplicationService) DeleteApplication(applicationID int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.applicationRepo.Delete(ctx, applicationID)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
