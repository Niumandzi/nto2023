package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type BookingRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewBookingRepository(db *sql.DB, logger logging.Logger) BookingRepository {
	return BookingRepository{
		db:     db,
		logger: logger,
	}
}

func (b BookingRepository) Create(ctx context.Context, booking model.Booking) (int, error) {
	var bookingId int64

	if len(booking.PartIDs) == 0 && booking.FacilityID == 0 {
		err := errors.New("no booking facilityID nor partIDs provided")
		b.logger.Logger.Errorf("error %v", err.Error())
		return 0, err
	} else if len(booking.PartIDs) > 0 && booking.FacilityID != 0 {
		err := errors.New("expected facilityID or partIDs got both instead")
		b.logger.Logger.Errorf("error %v", err.Error())
		return 0, err
	}

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		b.logger.Logger.Errorf("error %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	if len(booking.PartIDs) > 0 {
		res, err := tx.ExecContext(ctx, `INSERT INTO booking (description, create_date, start_date, end_date, event_id, facility_id) 
												VALUES ($1, $2, $3, $4, $5, $6);`, booking.Description, booking.CreateDate, booking.StartDate, booking.EndDate, booking.EventID, 0)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return 0, err
		}

		bookingId, err := res.LastInsertId()

		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return 0, err
		}
		for _, partID := range booking.PartIDs {
			_, err = tx.ExecContext(ctx, `INSERT INTO booking_part (booking_id, part_id) VALUES ($1, $2);`, bookingId, partID)
			if err != nil {
				b.logger.Errorf("error: %v", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	} else {
		res, err := tx.ExecContext(ctx, `INSERT INTO booking (description, create_date, start_date, end_date, event_id, facility_id) 
												VALUES ($1, $2, $3, $4, $5, $6);`, booking.Description, booking.CreateDate, booking.StartDate, booking.EndDate, booking.EventID, booking.FacilityID)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return 0, err
		}

		bookingId, err = res.LastInsertId()

		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(bookingId), nil
}

func (b BookingRepository) Get(ctx context.Context, startDate string, endDate string, eventId int, categoryName string) ([]model.BookingWithFacility, error) {
	var bookings []model.BookingWithFacility
	kwargs := make(map[string]interface{})
	args := make([]interface{}, 0, 3)

	query := `SELECT booking.id, 
				booking.description, 
				booking.create_date, 
				booking.start_date, 
				booking.end_date,
				booking.event_id, 
				booking.facility_id,
				facility.name, 
				facility.have_parts, 
				facility.is_active
			FROM
			    booking
			INNER JOIN 
			        events ON booking.event_id = events.id
			INNER JOIN 
			        details ON events.details_id = details.id
			INNER JOIN
			    facility ON booking.facility_id = facility.id`

	var dateQuery string
	if startDate != "" && endDate != "" {
		dateQuery = "(booking.end_date < ? OR booking.start_date > ?)"
	}
	if eventId != 0 {
		kwargs["booking.facility_id"] = eventId
	}
	if categoryName != "" {
		kwargs["details.category"] = categoryName
	}

	var rows *sql.Rows
	var err error

	length := len(kwargs)
	if length == 0 && dateQuery == "" {
		query += `;`
		rows, err = b.db.QueryContext(ctx, query)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			return []model.BookingWithFacility{}, nil
		}
	} else {
		query += ` WHERE `
		if dateQuery != "" {
			query += dateQuery + ` AND `
		}
		i := 0
		for key, val := range kwargs {
			if i == length-1 {
				query += fmt.Sprintf("%v = ?;", key)
				args = append(args, val)
			} else {
				query += fmt.Sprintf("%v ? AND ", key)
				args = append(args, val)
			}
			i++
		}
		rows, err = b.db.QueryContext(ctx, query, args...)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			return []model.BookingWithFacility{}, nil
		}
	}

	defer rows.Close()

	for rows.Next() {
		var booking model.BookingWithFacility
		var parts []model.Parts

		err = rows.Scan(&booking.ID,
			&booking.Description,
			&booking.CreateDate,
			&booking.EndDate,
			&booking.Event.ID,
			&booking.Facility.ID,
			&booking.Facility.Name,
			&booking.Facility.HaveParts)

		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			return []model.BookingWithFacility{}, err
		}

		partRows, err := b.db.QueryContext(ctx, `SELECT part.id, part.name, part.facility_id FROM booking_part 
    													INNER JOIN 
    															part ON booking_part.part_id = part.id 
                                            			WHERE booking_part.booking_id = ?`,
			booking.ID)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			return []model.BookingWithFacility{}, err
		}

		defer partRows.Close()

		for partRows.Next() {
			var part model.Parts
			err = partRows.Scan(&part.ID, part.Name, part.FacilityID)

			parts = append(parts, part)
		}

		booking.Parts = parts
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (b BookingRepository) Update(ctx context.Context, bookingUpd model.Booking) error {
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE booking
											SET 
												description = ?,
												create_date = ?,
												start_date = ?,
												end_date = ?,
												event_id = ?,
												facility_id = ?
											WHERE
											    id = ?`, bookingUpd.Description, bookingUpd.CreateDate,
		bookingUpd.EndDate, bookingUpd.EventID, bookingUpd.FacilityID,
		bookingUpd.ID)
	rowCount, err := res.RowsAffected()
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowCount != 1 {
		err = errcode.NewRowCountError("booking repo update", int(rowCount))
		b.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	for _, bookingPartID := range bookingUpd.PartIDs {
		res, err = tx.ExecContext(ctx, `UPDATE 
    											booking_part
											SET
											    part_id = ?
											WHERE booking_id = ?`, bookingPartID, bookingUpd.ID)
		if err != nil {
			b.logger.Errorf("error: %v", err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	return nil
}

func (b BookingRepository) Delete(ctx context.Context, bookingId int) error {
	_, err := b.db.ExecContext(ctx, `DELETE FROM booking WHERE booking.id = ?`, bookingId)
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		return err
	}
	return nil
}
