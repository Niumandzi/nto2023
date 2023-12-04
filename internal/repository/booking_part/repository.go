package booking_part

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type BookingPartRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewBookingPartRepository(db *sql.DB, logger logging.Logger) BookingPartRepository {
	return BookingPartRepository{
		db:     db,
		logger: logger,
	}
}

func (b BookingPartRepository) Create(ctx context.Context, bookingId int, partId int) error {
	_, err := b.db.ExecContext(ctx, `INSERT INTO booking_part (booking_id, part_id) VALUES (?, ?)`, bookingId, partId)
	if err != nil {
		b.logger.Error(err.Error())
		return err
	}
	return nil
}

func (b BookingPartRepository) GetParts(ctx context.Context, bookingId int) ([]int, error) {
	var partIDs []int

	rows, err := b.db.QueryContext(ctx, `SELECT part_id FROM booking_part WHERE booking_id = ?`, bookingId)
	if err != nil {
		b.logger.Error(err.Error())
		return []int{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var partID int

		err := rows.Scan(&partID)
		if err != nil {
			b.logger.Error(err.Error())
			return []int{}, err
		}

		partIDs = append(partIDs, partID)
	}

	return partIDs, nil
}

func (b BookingPartRepository) Delete(ctx context.Context, bookingId int, partId int) error {
	_, err := b.db.ExecContext(ctx, `DELETE FROM booking_part WHERE booking_id = ? AND part_id = ?`, bookingId, partId)
	if err != nil {
		b.logger.Error(err.Error())
		return err
	}
	return nil
}
