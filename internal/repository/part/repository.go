package part

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/pkg/logging"
	"strings"
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

func (p PartRepository) Create(ctx context.Context, facilityID int, partNames []string) (int, error) {
	println(facilityID, partNames)
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO part (name, facility_id) VALUES ")

	var args []interface{}
	for i, name := range partNames {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString("(?, ?)")
		args = append(args, name, facilityID)
	}

	res, err := tx.ExecContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

func (p PartRepository) Update(ctx context.Context, update map[int]string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		return err
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`UPDATE part SET name = CASE id `)
	var args []interface{}
	var i int = 1
	for id, name := range update {
		queryBuilder.WriteString(fmt.Sprintf("WHEN $%d THEN $%d ", i, i+1))
		args = append(args, id, name)
		i += 2
	}
	queryBuilder.WriteString("END WHERE id IN (")
	first := true
	j := 1
	for id := range update {
		if !first {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("$%d", i+j))
		args = append(args, id)
		j++
		first = false
	}
	queryBuilder.WriteString(");")

	res, err := tx.ExecContext(ctx, queryBuilder.String(), args...)
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

	if rowsCount != int64(len(update)) {
		err = errors.NewRowCountError("expected to update parts", len(update))
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

func (p PartRepository) Delete(ctx context.Context, delete map[int]bool) error {
	fmt.Print(delete)
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		return err
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`UPDATE part SET is_active = CASE id `)
	var args []interface{}
	var i int = 1
	for id, isActive := range delete {
		queryBuilder.WriteString(fmt.Sprintf("WHEN $%d THEN $%d ", i, i+1))
		args = append(args, id, isActive)
		i += 2
	}
	queryBuilder.WriteString("END WHERE id IN (")
	first := true
	j := 1
	for id := range delete {
		if !first {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("$%d", i+j))
		args = append(args, id)
		j++
		first = false
	}
	queryBuilder.WriteString(");")

	res, err := tx.ExecContext(ctx, queryBuilder.String(), args...)
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

	if rowsCount != int64(len(delete)) {
		err = errors.NewRowCountError("expected to delete/update parts", len(delete))
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
