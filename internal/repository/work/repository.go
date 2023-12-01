package work

import (
	"context"
	"database/sql"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type WorkTypeRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewWorkTypeRepository(db *sql.DB, logger logging.Logger) WorkTypeRepository {
	return WorkTypeRepository{
		db:     db,
		logger: logger,
	}
}

func (w WorkTypeRepository) Create(ctx context.Context, name string) (int, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO work_type (name) VALUES ($1);`, name)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (w WorkTypeRepository) Get(ctx context.Context, categoryName string, facilityID int, status string, isActive bool) ([]model.WorkType, error) {
	args := make([]interface{}, 0, 5)
	var workTypes []model.WorkType

	baseQuery := `SELECT work_type.id,
                         work_type.name
			FROM work_type
			LEFT JOIN application ON application.work_type_id = work_type.id
			LEFT JOIN facility ON application.facility_id = facility.id
			LEFT JOIN events ON application.event_id = events.id
			LEFT JOIN details ON events.details_id = details.id
			WHERE work_type.is_active = ?`

	args = append(args, isActive)

	if categoryName != "" {
		baseQuery += ` AND details.category = ?`
		args = append(args, categoryName)
	}
	if facilityID != 0 {
		baseQuery += ` AND facility.id = ?`
		args = append(args, facilityID)
	}
	if status != "" {
		baseQuery += ` AND application.status = ?`
		args = append(args, status)
	}

	rows, err := w.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var workType model.WorkType

		err = rows.Scan(&workType.ID, &workType.Name)
		if err != nil {
			w.logger.Errorf("error: %v", err.Error())
			return nil, err
		}

		workTypes = append(workTypes, workType)
	}

	return workTypes, nil
}

func (w WorkTypeRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE work_type SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("work type name update", int(rowsCount))
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (w WorkTypeRepository) Delete(ctx context.Context, workTypeId int, isActive bool) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE work_type SET is_active = $1 WHERE work_type.id = $2;`, isActive, workTypeId)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("work type delete", int(rowsCount))
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
