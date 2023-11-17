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

func (s EventService) DeleteType(categoryName string, typeName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.DeleteType(ctx, categoryName, typeName)
	if err != nil {
		return err
	}
	return nil
}
