package registration

import (
	"context"
	"database/sql"
	"errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type RegistrationRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewRegistrationRepository(db *sql.DB, logger logging.Logger) RegistrationRepository {
	return RegistrationRepository{
		db:     db,
		logger: logger,
	}
}

func (r RegistrationRepository) Create(ctx context.Context, registration model.Registration) (int, error) {
	var registrationID int64

	if len(registration.PartIDs) == 0 && registration.FacilityID == 0 || registration.FacilityID == 0 {
		err := errors.New("no booking facilityID no partIDs provided")
		r.logger.Logger.Error("error ", err.Error())
		return 0, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Logger.Error("error ", err.Error())
		tx.Rollback()
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO registration (name, start_date, number_of_days, facility_id, mug_type_id, teacher_id) 
												VALUES ($1, $2, $3, $4, $5, $6);`,
		registration.Name,
		registration.StartDate,
		registration.NumberOfDays,
		registration.FacilityID,
		registration.MugTypeID,
		registration.TeacherID)
	if err != nil {
		r.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	registrationID, err = res.LastInsertId()
	err = tx.Commit()
	if err != nil {
		r.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	if len(registration.Schedule) > 0 {
		if err != nil {
			r.logger.Error("error: ", err.Error())
			tx.Rollback()
			return 0, err
		}
		for _, schedule := range registration.Schedule {
			_, err = tx.ExecContext(ctx, `INSERT INTO schedule (registration_id, name, start_time, end_time) 
													VALUES ($1, $2, $3, $4);`,
				registrationID,
				schedule.Name,
				schedule.StartTime,
				schedule.EndTime)
			if err != nil {
				r.logger.Error("error: ", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	if len(registration.PartIDs) > 0 {
		if err != nil {
			r.logger.Error("error: ", err.Error())
			tx.Rollback()
			return 0, err
		}
		for _, partID := range registration.PartIDs {
			_, err = tx.ExecContext(ctx, `INSERT INTO registration_part (registration_id, part_id) VALUES ($1, $2);`, registrationID, partID)
			if err != nil {
				r.logger.Error("error: ", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(registrationID), nil
}

func (r RegistrationRepository) Get(ctx context.Context, facilityID int, mugID int, teacherID int) ([]model.RegistrationWithDetails, error) {
	return nil, nil
}
