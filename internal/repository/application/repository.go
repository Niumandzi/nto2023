package application

import (
	"context"
	"database/sql"
	"errors"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type ApplicationRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewApplicationRepository(db *sql.DB, logger logging.Logger) ApplicationRepository {
	return ApplicationRepository{
		db:     db,
		logger: logger,
	}
}

func (a ApplicationRepository) Create(ctx context.Context, application model.Application) (int, error) {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO application 
											(description, 
											 created_at, 
											 due, 
											 status, 
											 work_type_id, 
											 event_id, 
											 facility_id) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		application.Description,
		application.CreateDate,
		application.DueDate,
		application.Status,
		application.WorkTypeId,
		application.EventId,
		application.FacilityId)
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

// аналогично как и в events repo, делаем один метод на get по workType и status
// TODO: сделать базовое query и добавлять новое string по наличию поля,
func (a ApplicationRepository) Get(ctx context.Context, categoryName string, workType string, status string) ([]model.ApplicationWithDetails, error) {
	args := make([]interface{}, 0, 2)
	var query string
	var applications []model.ApplicationWithDetails

	baseQuery := `SELECT application.id,
		application.description,
		application.created_at,
		application.due,
		application.status,
		work_type.id,
		work_type.name,
		events.id,
        events.name,
		events.description,
		events.date,
        details.id,
        details.type_name,
        details.category,
		facility.id,
		facility.name
	FROM application
	INNER JOIN work_type ON application.work_type_id = work_type.id
	INNER JOIN facility ON application.facility_id = facility.id
	INNER JOIN events ON application.event_id = events.id
	INNER JOIN details ON events.details_id = details.id`

	if (categoryName == "") && (workType == "") && (status == "") {
		err := errors.New("no categoryName, workType and status provided")
		a.logger.Errorf(err.Error())
		return applications, err
	} else if (categoryName == "") && (workType == "") && (status != "") {
		baseQuery += `WHERE application.status = $1;`
		args = append(args, status)
	} else if (categoryName == "") && (workType != "") && (status != "") {
		baseQuery += `WHERE work_type.name = $1 AND application.status = $2;`
		args = append(args, workType, status)
	} else if (categoryName != "") && (workType == "") && (status == "") {
		baseQuery += `WHERE details.category = $1;`
		args = append(args, categoryName)
	} else if (categoryName != "") && (workType != "") && (status == "") {
		baseQuery += `WHERE details.category = $1 AND work_type.name = $2;`
		args = append(args, categoryName, workType)
	} else if (categoryName != "") && (workType == "") && (status != "") {
		baseQuery += `WHERE details.category = $1 AND application.status = $2;`
		args = append(args, categoryName, status)
	} else if (categoryName != "") && (workType != "") && (status != "") {
		baseQuery += `WHERE details.category = $1 AND work_type.name = $2 AND application.status = $3;`
		args = append(args, categoryName, workType, status)
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		a.logger.Error(err.Error())
		return applications, err
	}

	defer rows.Close()

	for rows.Next() {
		var application model.ApplicationWithDetails

		err = rows.Scan(&application.ID,
			&application.Description,
			&application.CreateDate,
			&application.DueDate,
			&application.Status,
			&application.WorkType.ID,
			&application.WorkType.Name,
			&application.Event.ID,
			&application.Event.Name,
			&application.Event.Description,
			&application.Event.Date,
			&application.Event.Details.ID,
			&application.Event.Details.Category,
			&application.Event.Details.TypeName,
			&application.Facility.ID,
			&application.Facility.Name,
		)

		if err != nil {
			a.logger.Errorf("error: %v", err.Error())
			return []model.ApplicationWithDetails{}, err
		}

		applications = append(applications, application)
	}

	return applications, nil
}

func (a ApplicationRepository) Update(ctx context.Context, applicationUpd model.Application) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE application 
											SET description = $1, 
											    created_at = $2, 
											    due = $3, 
											    status = $4, 
											    work_type_id = $5, 
											    event_id = $6, 
											    facility_id = $7 
											WHERE id = $8;`,
		applicationUpd.Description,
		applicationUpd.CreateDate,
		applicationUpd.DueDate,
		applicationUpd.Status,
		applicationUpd.WorkTypeId,
		applicationUpd.EventId,
		applicationUpd.FacilityId,
		applicationUpd.ID)

	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("application update", int(rowsCount))
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (a ApplicationRepository) Delete(ctx context.Context, applicationId int) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM application WHERE id = $1;`, applicationId)
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("application delete", int(rowsCount))
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		a.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
