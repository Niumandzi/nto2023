package event

import "context"

func (s EventService) DeleteEvent(eventId int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.eventRepo.Delete(ctx, eventId)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
