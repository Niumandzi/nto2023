package event

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type EventRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewEventRepository(db *sql.DB, logger logging.Logger) EventRepository {
	return EventRepository{
		db:     db,
		logger: logger,
	}
}

func (s EventRepository) Create(ctx context.Context, event model.EventWithCategoryAndType) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return 0, nil
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO events (name, description, date, details_id) VALUES ($1, $2, $3, $4, $5);`, event.Name, event.Description, event.Date, event.Category.ID)
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

// Get объединяем два запроса в один, выбор запроса зависит от eventType.
// Он может быть либо по event_type_id или по event type, либо по category.
func (s EventRepository) Get(ctx context.Context, eventCategory string, eventType string) ([]model.EventWithCategoryAndType, error) {
	var query string
	var events []model.EventWithCategoryAndType

	switch eventType {
	//запрос только по по категории
	case "":
		query = `SELECT events.id,
       					events.name,
						events.description,
						events.date,
						categories_types.id,
						categories_types.category,
						event_types.id,
						event_types.type_name
				   FROM events
				   INNER JOIN categories_types ON events.details_id = categories_types.id
				   INNER JOIN event_types ON categories_types.type_id = event_types.id
				   WHERE categories_types.category = $1;`
		break
	//запрос категории + тип
	default:
		query = `SELECT events.id, 
						events.name, 
						events.description, 
						events.date, 
						categories_types.id,  
						categories_types.category, 
						event_types.id, 
						event_types.type_name
				FROM events
				INNER JOIN categories_types ON events.details_id = categories_types.id
				INNER JOIN event_types ON categories_types.type_id = event_types.id
				WHERE categories_types.category = $1 AND event_types.type_name = $2;`
	}

	rows, err := s.db.QueryContext(ctx, query, eventCategory, eventType)
	if err != nil {
		s.logger.Error(err.Error())
		return []model.EventWithCategoryAndType{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var event model.EventWithCategoryAndType

		err = rows.Scan(&event.ID,
			&event.Name,
			&event.Description,
			&event.Date,
			&event.Category.ID,
			&event.Category.Category,
			&event.Category.EventType.ID,
			&event.Category.EventType.TypeName,
		)

		if err != nil {
			s.logger.Fatalf("error: %v", err.Error())
			return []model.EventWithCategoryAndType{}, nil
		}

		events = append(events, event)
	}

	return events, nil
}

// Update обновляет только type_id, name, date, description.
// Обновление category у event type и изменение самой category - это отдельные методы в отдельном репо.
func (s EventRepository) Update(ctx context.Context, eventUpd model.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return nil
	}

	res, err := tx.ExecContext(ctx, `UPDATE events SET name = $1, description = $2, date = $3, details_id = $4 WHERE id = $5;`, eventUpd.Name, eventUpd.Description, eventUpd.Date, eventUpd.DetailsID, eventUpd.ID)
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
		err = errors.NewRowCountError("event update", int(rowsCount))
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

func (s EventRepository) Delete(ctx context.Context, eventId int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Fatalf("error: %v", err.Error())
		return nil
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM events WHERE id = $5;`, eventId)
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
		err = errors.NewRowCountError("event delete", int(rowsCount))
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
