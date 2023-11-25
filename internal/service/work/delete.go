package work

import "context"

func (s WorkTypeService) DeleteWorkType(id int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.workTypeRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
