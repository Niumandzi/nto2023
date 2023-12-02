package part

import "context"

func (s PartService) DeleteBooking(partIds []int, isActive bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := s.partRepo.Delete(ctx, partIds, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
