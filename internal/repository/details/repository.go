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

	rows, err := s.db.QueryContext(ctx, `SELECT details.id, details.type_name, details.category FROM details WHERE details.category = $1`, categoryName)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return []model.Details{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var detail model.Details

		err = rows.Scan(&detail.Id,
			&detail.TypeName,
			&detail.Category)
		if err != nil {
			s.logger.Errorf("error: %v", err.Error())
			return []model.Details{}, err
		}

		details = append(details, detail)
	}

	return details, nil
}

func (s DetailsRepository) GetId(ctx context.Context, categoryName string, typeName string) (int, error) {
	var id int

	row := s.db.QueryRowContext(ctx, `SELECT details.id FROM details WHERE details.category = $1 AND details.type_name = $2`, categoryName, typeName)

	err := row.Scan(&id)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	return id, nil
}

func (s DetailsRepository) UpdateTypeName(ctx context.Context, detailsId int, typeName string) error {
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
		s.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (s DetailsRepository) DeleteType(ctx context.Context, detailsId int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM details WHERE details.id = $1;`, detailsId)
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
