package part

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
	"strconv"
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

func (p PartRepository) Create(ctx context.Context, part model.Part) (int, error) {
	return 0, nil
}

func (p PartRepository) Update(ctx context.Context, updates map[int]string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Errorf("error: %v", err.Error())
		return err
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`UPDATE part SET name = CASE id `)
	var args []interface{}
	var i int = 1
	for id, name := range updates {
		queryBuilder.WriteString(fmt.Sprintf("WHEN $%d THEN $%d ", i, i+1))
		args = append(args, id, name)
		i += 2
	}
	queryBuilder.WriteString("END WHERE id IN (")
	first := true
	j := 1
	for id := range updates {
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

	if rowsCount != int64(len(updates)) {
		err = errors.NewRowCountError("expected to update parts", len(updates))
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

func (p PartRepository) Delete(ctx context.Context, partIds []int, isActive bool) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error("error: %v", err.Error())
		return err
	}

	query := `UPDATE part SET is_active = $1 WHERE part.id IN (`
	params := []interface{}{isActive}

	for i, partID := range partIds {
		query += "$" + strconv.Itoa(i+2)
		params = append(params, partID)

		if i < len(partIds)-1 {
			query += ", "
		}
	}
	query += ");"

	res, err := tx.ExecContext(ctx, query, params...)
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
	if rowsCount != int64(len(partIds)) {
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
