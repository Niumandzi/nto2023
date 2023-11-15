package event

import (
	"context"
	"database/sql"
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

func (s EventRepository) Create(ctx context.Context, contact model.Event) (int, error) {

	return 0, err
}

// Get объединяем два запроса в один, выбор запроса зависит от eventArgument.
// Он может быть либо по event_type_id или по event type, либо по category.
func (s EventRepository) Get(ctx context.Context, eventArgument string) ([]model.EventWithCategoryAndType, error) {
	query := "" //запрос по event_type_id или по event type
	var args []interface{}

	if eventArgument != "entertainment" && eventArgument != "enlightenment" && eventArgument != "education" {
		query = "" //запрос по категории
		args = append(args, eventArgument)
	}

	return nil, err
}

// Update обновляет только type_id, name, date, description.
// Обновление category у event type и изменение самой category - это отдельные методы в отдельном репо.
func (s EventRepository) Update(ctx context.Context, eventInput model.Event) error {

	return nil
}

func (s EventRepository) Delete(ctx context.Context, eventId int) error {

	return nil
}
