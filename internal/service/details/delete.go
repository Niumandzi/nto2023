package event

import "context"

func (s DetailsService) DeleteRestoreType(detailsID int, isActive bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.DeleteRestoreType(ctx, detailsID, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
