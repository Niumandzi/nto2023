package application

import (
	"context"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (s ApplicationService) UpdateApplication(applicationUpd model.Application) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&applicationUpd,
		validation.Field(&applicationUpd.Description, validation.Required),
		validation.Field(&applicationUpd.CreateDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&applicationUpd.DueDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&applicationUpd.Status, validation.Required, validation.In("created", "todo", "done")),
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
