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

func (s EventRepository) Create(ctx context.Context, event model.Event) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO events (name, description, date, details_id) VALUES ($1, $2, $3, $4);`, event.Name, event.Description, event.Date, event.DetailsID)
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

// Get объединяем два запроса в один, выбор запроса зависит от eventType.
// Он может быть либо по event_type_id или по event type, либо по category.
func (s EventRepository) Get(ctx context.Context, categoryName string, detailsID int) ([]model.EventWithDetails, error) {
	var query string
	var args []interface{}
	var events []model.EventWithDetails

	switch detailsID {
	case -1:
		query = `SELECT events.id,
                     events.name,
                     events.description,
                     events.date,
                     details.id,
                     details.type_name,
                     details.category
              FROM events
              INNER JOIN details ON events.details_id = details.id
              WHERE details.category = $1;`
		args = append(args, categoryName)

	default:
		query = `SELECT events.id, 
                     events.name, 
                     events.description, 
                     events.date, 
                     details.id,
                     details.type_name,
                     details.category
             FROM events
             INNER JOIN details ON events.details_id = details.id
             WHERE details.id = $1;`
		args = append(args, detailsID)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Error(err.Error())
		return []model.EventWithDetails{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var event model.EventWithDetails

		err = rows.Scan(&event.ID,
			&event.Name,
			&event.Description,
			&event.Date,
			&event.Details.ID,
			&event.Details.TypeName,
			&event.Details.Category,
		)

		if err != nil {
			s.logger.Errorf("error: %v", err.Error())
			return []model.EventWithDetails{}, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (s EventRepository) Update(ctx context.Context, eventUpd model.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE events SET name = $1, description = $2, date = $3, details_id = $4 WHERE id = $5;`, eventUpd.Name, eventUpd.Description, eventUpd.Date, eventUpd.DetailsID, eventUpd.ID)
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
		err = errors.NewRowCountError("event update", int(rowsCount))
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

func (s EventRepository) Delete(ctx context.Context, eventId int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM events WHERE id = $1;`, eventId)
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
		err = errors.NewRowCountError("event delete", int(rowsCount))
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
