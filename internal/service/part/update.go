package part

import (
	"context"
)

func (s PartService) UpdatePart(update map[int]string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.partRepo.Update(ctx, update)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
