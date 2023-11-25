package application

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/internal/validations"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationService) UpdateApplication(applicationUpd model.Application) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&applicationUpd,
		validation.Field(&applicationUpd.ID, validation.Required, validation.Min(1)),
		validation.Field(&applicationUpd.Description, validation.Required),
		validation.Field(&applicationUpd.CreateDate, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&applicationUpd.DueDate, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&applicationUpd.Description, validation.Required, validation.In("todo", "done")),
		validation.Field(&applicationUpd.EventId, validation.Required, validation.Min(1).Error("Не выбрано мероприятие")),
		validation.Field(&applicationUpd.WorkTypeId, validation.Required, validation.Min(1).Error("Не выбрано помещение")),
		validation.Field(&applicationUpd.FacilityId, validation.Required, validation.Min(1).Error("Не выбран тип работ")),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	err = s.applicationRepo.Update(ctx, applicationUpd)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
