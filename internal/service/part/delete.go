package part

import "context"

func (s PartService) DeletePart(delete map[int]bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.partRepo.Delete(ctx, delete)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
