package registration

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
	"time"
)

func (r RegistrationService) CreateRegistration(registration model.Registration) (int, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.contextTimeout)
	defer cancel()

	err := checkForDuplicateDays(registration.Schedule)
	if err != nil {
		r.logger.Error("error: %v", err.Error())
		return 0, err
	}

	err = validation.ValidateStruct(&registration,
		validation.Field(&registration.Name, validation.Required),
		validation.Field(&registration.StartDate, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&registration.NumberOfDays, validation.Required, validation.In(1, 2, 3)),
		validation.Field(&registration.FacilityID, validation.Required, validation.Min(1).Error("Не выбрано помещение")),
		validation.Field(&registration.MugTypeID, validation.Required, validation.Min(1).Error("Не выбран тип кружка")),
		validation.Field(&registration.TeacherID, validation.Required, validation.Min(1).Error("Не выбран преподаватель")),
		validation.Field(&registration.Schedule, validation.Required, validation.Each(validation.Required, validation.By(validateSchedule))),
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

func checkForDuplicateDays(schedule []model.Schedule) error {
	days := make(map[string]bool)
	for _, s := range schedule {
		if days[s.Day] {
			return errors.New("duplicate day found in schedule")
		}
		days[s.Day] = true
	}
	return nil
}

func validateSchedule(value interface{}) error {
	schedule, ok := value.(model.Schedule)
	if !ok {
		return errors.New("invalid schedule type")
	}

	return validation.ValidateStruct(&schedule,
		validation.Field(&schedule.StartTime, validation.Required, validation.Date("15:04")),
		validation.Field(&schedule.EndTime, validation.Required, validation.Date("15:04"), validation.By(func(endTime interface{}) error {
			st, _ := time.Parse("15:04", schedule.StartTime)
			et, _ := time.Parse("15:04", endTime.(string))
			if et.Before(st) {
				return errors.New("EndTime must be after StartTime")
			}
			return nil
		})),
	)
}
