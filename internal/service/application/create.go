package application

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/internal/validations"
	"github.com/niumandzi/nto2023/model"
	"golang.org/x/net/context"
)

func (s ApplicationService) CreateApplication(application model.Application) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&application,
		validation.Field(&application.Description, validation.Required),
		validation.Field(&application.CreateDate, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&application.DueDate, validation.Required, validation.By(validations.ValidateDate)),
		validation.Field(&application.Status, validation.Required, validation.In("created")),
		validation.Field(&application.EventId, validation.Required, validation.Min(1).Error("Не выбрано мероприятие")),
		validation.Field(&application.WorkTypeId, validation.Required, validation.Min(1).Error("Не выбрано помещение")),
		validation.Field(&application.FacilityId, validation.Required, validation.Min(1).Error("Не выбран тип работ")),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	applicationDB := model.Application{
		Description: application.Description,
		CreateDate:  application.CreateDate,
		DueDate:     application.DueDate,
		Status:      application.Status,
		EventId:     application.EventId,
		WorkTypeId:  application.WorkTypeId,
		FacilityId:  application.FacilityId,
	}

	id, err := s.applicationRepo.Create(ctx, applicationDB)
	if err != nil {
		return 0, err
	}

	return id, nil
}
