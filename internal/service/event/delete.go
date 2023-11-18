package event

import "context"

func (s EventService) DeleteEvent(eventId int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.eventRepo.Delete(ctx, eventId)
	if err != nil {
		return err
	}

	return nil
}

func (s EventService) DeleteType(detailsId int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.DeleteType(ctx, detailsId)
	if err != nil {
		return err
	}
	return nil
}
