package application

import "context"

func (s ApplicationService) DeleteApplication(applicationId int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.applicationRepo.Delete(ctx, applicationId)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
