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
		return err
	}

	err = s.eventRepo.Update(ctx, eventUpd)
	if err != nil {
		return err
	}

	return nil
}

func (s EventService) UpdateTypeName(eventCategory string, eventType string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	typeID, err := s.categoryTypeRepo.GetCategoryTypeID(ctx, eventCategory, eventType)
	if err != nil {
		return err
	}

	err = s.categoryTypeRepo.UpdateTypeName(ctx, typeID, eventCategory)
	if err != nil {
		return err
	}

	return nil
}
