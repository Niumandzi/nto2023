package part

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s PartService) UpdateBooking(partId int, name string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.partRepo.Update(ctx, partId, name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
