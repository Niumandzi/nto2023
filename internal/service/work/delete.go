package work

import "context"

func (s WorkTypeService) DeleteRestoreWorkType(workTypeID int, isActive bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.workTypeRepo.Delete(ctx, workTypeID, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
