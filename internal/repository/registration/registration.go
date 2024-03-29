package registration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
	"strconv"
	"strings"
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
	fmt.Print(registration)

	if len(registration.PartIDs) == 0 && registration.FacilityID == 0 {
		err := errors.New("no booking facilityID no partIDs provided")
		r.logger.Logger.Error("error ", err.Error())
		return 0, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Logger.Error("error ", err.Error())
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			r.logger.Logger.Error("Transaction rolled back due to error: ", err.Error())
		}
	}()

	res, err := tx.ExecContext(ctx, `INSERT INTO registration (name, start_date, number_of_days, facility_id, mug_type_id, teacher_id) 
												VALUES ($1, $2, $3, $4, $5, $6);`,
		registration.Name,
		registration.StartDate,
		registration.NumberOfDays,
		registration.FacilityID,
		registration.MugTypeID,
		registration.TeacherID)
	if err != nil {
		r.logger.Error("registration error: ", err.Error())
		return 0, err
	}

	registrationID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Print(registration.Schedule)
	if len(registration.Schedule) > 0 {
		for _, schedule := range registration.Schedule {
			_, err = tx.ExecContext(ctx, `INSERT INTO schedule (registration_id, day, start_time, end_time) 
													VALUES ($1, $2, $3, $4);`,
				registrationID,
				schedule.Day,
				schedule.StartTime,
				schedule.EndTime)
			if err != nil {
				r.logger.Error("schedule error: ", err.Error())
				return 0, err
			}
		}
	}

	if len(registration.PartIDs) > 0 {
		for _, partID := range registration.PartIDs {
			_, err = tx.ExecContext(ctx, `INSERT INTO registration_part (registration_id, part_id) VALUES ($1, $2);`, registrationID, partID)
			if err != nil {
				r.logger.Error("part error: ", err.Error())
				return 0, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		r.logger.Logger.Error("error: ", err.Error())
		return 0, err
	}

	return int(registrationID), nil
}

func (r RegistrationRepository) Get(ctx context.Context, facilityID int, mugID int, teacherID int) ([]model.RegistrationWithDetails, error) {
	var registrations []model.RegistrationWithDetails
	args := make([]interface{}, 0)

	query := `SELECT registration.id,
				registration.name,
				registration.start_date,
				registration.number_of_days,
				registration.facility_id,
				facility.name, 
				facility.have_parts,
				registration.mug_type_id,
				mug_type.name,
				registration.teacher_id,
				teacher.name,
				COALESCE(GROUP_CONCAT(schedule.id), '')  AS schedule_ids,
				COALESCE(GROUP_CONCAT(schedule.day), '')  AS schedule_days,
				COALESCE(GROUP_CONCAT(schedule.start_time), '')  AS schedule_start_times,
				COALESCE(GROUP_CONCAT(schedule.end_time), '')  AS schedule_end_times,
				COALESCE(GROUP_CONCAT(part.id), '')  AS part_ids,
				COALESCE(GROUP_CONCAT(part.name), '')  AS part_names
		FROM registration
		INNER JOIN facility ON registration.facility_id = facility.id
		INNER JOIN mug_type ON registration.mug_type_id = mug_type.id
		INNER JOIN teacher ON registration.teacher_id = teacher.id
		LEFT JOIN schedule ON registration.id = schedule.registration_id
		LEFT JOIN registration_part ON registration.id = registration_part.registration_id
		LEFT JOIN part ON registration_part.part_id = part.id GROUP BY registration.id;`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var registration model.RegistrationWithDetails
		var scheduleIDs, scheduleDays, scheduleStartTime, scheduleEndTime, partIDs, partNames string

		err = rows.Scan(
			&registration.ID,
			&registration.Name,
			&registration.StartDate,
			&registration.NumberOfDays,
			&registration.Facility.ID,
			&registration.Facility.Name,
			&registration.Facility.HaveParts,
			&registration.MugType.ID,
			&registration.MugType.Name,
			&registration.Teacher.ID,
			&registration.Teacher.Name,
			&scheduleIDs,
			&scheduleDays,
			&scheduleStartTime,
			&scheduleEndTime,
			&partIDs,
			&partNames,
		)

		if err != nil {
			r.logger.Errorf("error scanning registration: %v", err.Error())
			return nil, err
		}

		ids := strings.Split(scheduleIDs, ",")
		days := strings.Split(scheduleDays, ",")
		start := strings.Split(scheduleStartTime, ",")
		end := strings.Split(scheduleEndTime, ",")
		for i := 0; i < registration.NumberOfDays; i++ {
			var schedule model.Schedule
			schedule.ID, err = strconv.Atoi(ids[i])
			if err != nil {
				r.logger.Errorf("error converting schedule ID to integer: %v", err.Error())
				continue
			}
			schedule.Day = days[i]
			schedule.StartTime = start[i]
			schedule.EndTime = end[i]
			schedule.RegistrationID = registration.ID
			registration.Schedule = append(registration.Schedule, schedule)
		}

		if registration.Facility.HaveParts && partIDs != "" {
			ids = strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			for i, idStr := range ids {
				var part model.Part
				part.ID, err = strconv.Atoi(idStr)
				if err != nil {
					r.logger.Errorf("error converting part ID to integer: %v", err.Error())
					continue
				}
				part.Name = names[i]
				part.FacilityID = registration.Facility.ID
				registration.Parts = append(registration.Parts, part)
			}
		}
		registrations = append(registrations, registration)
	}

	return registrations, nil
}

func (r RegistrationRepository) Update(ctx context.Context, registrationUpd model.Registration) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE registration
											SET 
												name = ?,
												number_of_days = ?,
												start_date = ?,
												facility_id = ?,
												mug_type_id = ?,
												teacher_id = ?
											WHERE
											    id = ?`, registrationUpd.Name, registrationUpd.NumberOfDays, registrationUpd.StartDate, registrationUpd.FacilityID,
		registrationUpd.MugTypeID, registrationUpd.TeacherID, registrationUpd.ID)
	rowCount, err := res.RowsAffected()
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	if rowCount != 1 {
		err = errcode.NewRowCountError("registration repo update", int(rowCount))
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM schedule WHERE registration_id = ?`, registrationUpd.ID)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	for _, schedule := range registrationUpd.Schedule {
		_, err = tx.ExecContext(ctx, `INSERT INTO schedule (day, start_time, end_time, registration_id) VALUES (?, ?, ?, ?)`, schedule.Day, schedule.StartTime, schedule.EndTime, registrationUpd.ID)
		if err != nil {
			r.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM registration_part WHERE registration_id = ?`, registrationUpd.ID)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	for _, partID := range registrationUpd.PartIDs {
		_, err = tx.ExecContext(ctx, `INSERT INTO registration_part (registration_id, part_id) VALUES (?, ?)`, registrationUpd.ID, partID)
		if err != nil {
			r.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (r RegistrationRepository) Delete(ctx context.Context, registrationID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM registration WHERE registration.id = ?`, registrationID)
	if err != nil {
		r.logger.Errorf("error: %v", err.Error())
		return err
	}
	return nil
}
