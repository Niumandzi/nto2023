package facility

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type FacilityRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewFacilityRepository(db *sql.DB, logger logging.Logger) FacilityRepository {
	return FacilityRepository{
		db:     db,
		logger: logger,
	}
}

func (w FacilityRepository) Create(ctx context.Context, name string) (int, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO facility (name) VALUES ($1);`, name)
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

func (w FacilityRepository) Get(ctx context.Context) ([]model.Facility, error) {
	var facilities []model.Facility

	rows, err := w.db.QueryContext(ctx, `SELECT facility.id, facility.name FROM facility`)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return []model.Facility{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var facility model.Facility

		err = rows.Scan(&facility.ID,
			&facility.Name)
		if err != nil {
			w.logger.Errorf("error: %v", err.Error())
			return []model.Facility{}, err
		}

		facilities = append(facilities, facility)
	}

	return facilities, nil
}

func (w FacilityRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET name = $1 WHERE id = $2;`, nameUpd, idOld)
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
		err = errors.NewRowCountError("facility name update", int(rowsCount))
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

func (w FacilityRepository) Delete(ctx context.Context, id int) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM facility WHERE facility.id = $1;`, id)
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
		err = errors.NewRowCountError("facility delete", int(rowsCount))
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
