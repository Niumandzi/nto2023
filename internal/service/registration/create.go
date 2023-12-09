package registration

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (r RegistrationService) CreateRegistration(registration model.Registration) (int, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.contextTimeout)
	defer cancel()

	err := validation.ValidateStruct(&registration,
		validation.Field(&registration.Name, validation.Required),
		validation.Field(&registration.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&registration.NumberOfDays, validation.Required, validation.In(1, 2, 3)),
		validation.Field(&registration.FacilityID, validation.Required, validation.Min(1).Error("Не выбрано помещение")),
		validation.Field(&registration.MugTypeID, validation.Required, validation.Min(1).Error("Не выбран тип кружка")),
		validation.Field(&registration.TeacherID, validation.Required, validation.Min(1).Error("Не выбран преподаватель")),
		validation.Field(&registration.Schedule, validation.Required),
	)
	if err != nil {
		r.logger.Error("error: %v", err.Error())
		return 0, err
	}

	id, err := r.registrationRepo.Create(ctx, registration)
	if err != nil {
		return 0, err
	}

	return id, nil
}
