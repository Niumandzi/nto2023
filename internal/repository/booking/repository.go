package booking

import (
	"context"
	"database/sql"
	"errors"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
	"strconv"
	"strings"
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
		err := errors.New("no booking facilityID no partIDs provided")
		b.logger.Logger.Errorf("error %v", err.Error())
		return 0, err
	}

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		b.logger.Logger.Error("error ", err.Error())
		tx.Rollback()
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO booking (description, create_date, start_date, end_date, event_id, facility_id) 
												VALUES ($1, $2, $3, $4, $5, $6);`, booking.Description, booking.CreateDate, booking.StartDate, booking.EndDate, booking.EventID, booking.FacilityID)
	if err != nil {
		b.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	bookingId, err = res.LastInsertId()
	if len(booking.PartIDs) > 0 {
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
	args := make([]interface{}, 0)

	query := `SELECT booking.id, 
                booking.description, 
                booking.create_date, 
                booking.start_date, 
                booking.end_date,
                booking.facility_id,
                facility.name, 
                facility.have_parts,
                booking.event_id, 
                events.name,
                events.date,
                events.description,
                details.category,
                details.type_name,
                COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
                COALESCE(GROUP_CONCAT(part.name), '') AS part_names
            FROM
                booking
            INNER JOIN 
                events ON booking.event_id = events.id
            INNER JOIN 
                details ON events.details_id = details.id
            INNER JOIN
                facility ON booking.facility_id = facility.id
            LEFT JOIN
                booking_part ON booking.id = booking_part.booking_id
            LEFT JOIN
                part ON booking_part.part_id = part.id`

	var whereClauses []string
	if startDate != "" && endDate != "" {
		whereClauses = append(whereClauses, "booking.start_date >= ? AND booking.end_date <= ?")
		args = append(args, startDate, endDate)
	}
	if eventId != 0 {
		whereClauses = append(whereClauses, "booking.event_id = ?")
		args = append(args, eventId)
	}
	if categoryName != "" {
		whereClauses = append(whereClauses, "details.category = ?")
		args = append(args, categoryName)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " GROUP BY booking.id;"

	rows, err := b.db.QueryContext(ctx, query, args...)
	if err != nil {
		b.logger.Errorf("error: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking model.BookingWithFacility
		var partIds, partNames string

		err = rows.Scan(
			&booking.ID,
			&booking.Description,
			&booking.CreateDate,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Facility.ID,
			&booking.Facility.Name,
			&booking.Facility.HaveParts,
			&booking.Event.ID,
			&booking.Event.Name,
			&booking.Event.Date,
			&booking.Event.Description,
			&booking.Event.Details.Category,
			&booking.Event.Details.TypeName,
			&partIds,
			&partNames,
		)
		if err != nil {
			b.logger.Errorf("error scanning booking: %v", err.Error())
			return nil, err
		}

		if booking.Facility.HaveParts && partIds != "" {
			ids := strings.Split(partIds, ",")
			names := strings.Split(partNames, ",")
			for i, idStr := range ids {
				var part model.Part
				part.ID, err = strconv.Atoi(idStr)
				if err != nil {
					b.logger.Errorf("error converting part ID to integer: %v", err.Error())
					continue
				}
				part.Name = names[i]
				part.FacilityID = booking.Facility.ID
				booking.Parts = append(booking.Parts, part)
			}
		}

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
