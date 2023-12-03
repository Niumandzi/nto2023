package details

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type DetailsRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewDetailsRepository(db *sql.DB, logger logging.Logger) *DetailsRepository {
	return &DetailsRepository{
		db:     db,
		logger: logger,
	}
}

func (s DetailsRepository) Create(ctx context.Context, categoryName string, typeName string) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO details (type_name, category) VALUES ($1, $2);`, typeName, categoryName)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (s DetailsRepository) Get(ctx context.Context, categoryName string) ([]model.Details, error) {
	var details []model.Details

	rows, err := s.db.QueryContext(ctx, `SELECT details.id, details.type_name, details.category, details.is_active FROM details WHERE details.category = $1`, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.Details{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var detail model.Details

		err = rows.Scan(
			&detail.ID,
			&detail.TypeName,
			&detail.Category,
			&detail.IsActive,
		)
		if err != nil {
			s.logger.Errorf("error: %v", err.Error())
			return []model.Details{}, err
		}

		details = append(details, detail)
	}

	return details, nil
}

func (s DetailsRepository) GetActive(ctx context.Context, categoryName string) ([]model.Details, error) {
	var details []model.Details

	rows, err := s.db.QueryContext(ctx, `SELECT details.id, details.type_name, details.category, details.is_active FROM details WHERE details.category = $1 and details.is_active = TRUE`, categoryName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return []model.Details{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var detail model.Details

		err = rows.Scan(
			&detail.ID,
			&detail.TypeName,
			&detail.Category,
			&detail.IsActive,
		)
		if err != nil {
			s.logger.Errorf("error: %v", err.Error())
			return []model.Details{}, err
		}

		details = append(details, detail)
	}

	return details, nil
}

func (s DetailsRepository) Update(ctx context.Context, detailsId int, typeName string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE details SET type_name = $1 WHERE id = $2;`, typeName, detailsId)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("details type name update", int(rowsCount))
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (s DetailsRepository) Delete(ctx context.Context, detailsId int, isActive bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE details SET is_active = $1 WHERE details.id = $2;`, isActive, detailsId)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("details delete", int(rowsCount))
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
