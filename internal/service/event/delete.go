package event

import "context"

func (s EventService) DeleteRestoreEvent(eventID int, isActive bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.eventRepo.Delete(ctx, eventID, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
