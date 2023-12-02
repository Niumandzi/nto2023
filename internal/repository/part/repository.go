package part

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type PartRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewPartRepository(db *sql.DB, logger logging.Logger) PartRepository {
	return PartRepository{
		db:     db,
		logger: logger,
	}
}

func (p PartRepository) Create(ctx context.Context, part model.Part) (int, error) {
	return 0, nil
}

func (p PartRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE part SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility name update", int(rowsCount))
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (p PartRepository) Delete(ctx context.Context, partId int, isActive bool) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE part SET is_active = $1 WHERE part.id = $2;`, isActive, partId)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("part delete", int(rowsCount))
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
