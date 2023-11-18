package event

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (s EventService) UpdateEvent(eventUpd model.Event) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&eventUpd, validation.Field(&eventUpd.ID, validation.Required, validation.Min(1)))
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	err = s.eventRepo.Update(ctx, eventUpd)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}

func (s EventService) UpdateTypeName(detailsId int, typeName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.UpdateTypeName(ctx, detailsId, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
