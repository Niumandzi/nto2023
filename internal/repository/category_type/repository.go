package category_type

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type CategoryTypeRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewCategoryTypeRepository(db *sql.DB, logger logging.Logger) *CategoryTypeRepository {
	return &CategoryTypeRepository{
		db:     db,
		logger: logger,
	}
}

func (s CategoryTypeRepository) CreateCategoryWithType(ctx context.Context, eventCategory string, eventType string) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, nil
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO event_types (type_name) VALUES ($1);`, eventType)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	eventTypeID, err := res.LastInsertId()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	res, err = tx.ExecContext(ctx, `INSERT INTO categories_types (type_id, category) VALUES ($1, $2)`, eventTypeID, eventCategory)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (s CategoryTypeRepository) GetCategoryTypes(ctx context.Context, eventCategory string) ([]model.EventType, error) {
	var eventTypes []model.EventType

	rows, err := s.db.QueryContext(ctx, `SELECT (event_types.id, event_types.type_name) FROM event_types INNER JOIN categories_types ON categories_types.type_id = event_types.id WHERE categories_types.category = $1`, eventCategory)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return []model.EventType{}, err
	}

	for rows.Next() {
		var eventType model.EventType

		err = rows.Scan(&eventType.ID,
			&eventType.TypeName)
		if err != nil {
			s.logger.Fatalf("error: %v", err.Error())
			return []model.EventType{}, err
		}

		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes, nil
}

func (s CategoryTypeRepository) GetCategoryTypeID(ctx context.Context, eventCategory string, eventType string) (int, error) {
	var id int

	row, err := s.db.QueryContext(ctx, `SELECT categories_types.id FROM categories_types INNER JOIN event_types ON event_types.id = categories_types.type_id WHERE category = $1 AND event_types.type_name = $2`, eventCategory, eventType)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, err
	}

	err = row.Scan(&id)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, err
	}

	return id, nil
}

func (s CategoryTypeRepository) UpdateTypeName(ctx context.Context, eventTypeID int, eventTypeUpd string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return nil
	}

	res, err := tx.ExecContext(ctx, `UPDATE event_types SET type_name = $1 WHERE id = $2;`, eventTypeUpd, eventTypeID)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("TypeName update", int(rowsCount))
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (s CategoryTypeRepository) DeleteType(ctx context.Context, eventType string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return nil
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM event_types
											WHERE type_name = $1;`, eventType)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("Type delete", int(rowsCount))
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
