package facility

import (
	"context"
	"database/sql"
	"github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
	"strconv"
	"strings"
)

type FacilityRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewFacilityRepository(db *sql.DB, logger logging.Logger) FacilityRepository {
	return FacilityRepository{
		db:     db,
		logger: logger,
	}
}

func (f FacilityRepository) Create(ctx context.Context, name string, parts []string) (int, error) {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	haveParts := len(parts) > 0
	res, err := tx.ExecContext(ctx, `INSERT INTO facility (name, have_parts) VALUES ($1, $2);`, name, haveParts)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	if haveParts {
		for _, partName := range parts {
			_, err := tx.ExecContext(ctx, `INSERT INTO part (name, facility_id) VALUES ($1, $2);`, partName, id)
			if err != nil {
				f.logger.Errorf("error: %v", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (f FacilityRepository) Get(ctx context.Context) ([]model.FacilityWithParts, error) {

	baseQuery := `SELECT facility.id, facility.name, facility.have_parts, facility.is_active,
                  COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
                  COALESCE(GROUP_CONCAT(part.name), '') AS part_names,
                  COALESCE(GROUP_CONCAT(part.is_active), '') AS part_is_active
           FROM facility
           LEFT JOIN part ON facility.id = part.facility_id
           LEFT JOIN application ON application.facility_id = facility.id GROUP BY facility.id`

	rows, err := f.db.QueryContext(ctx, baseQuery)
	if err != nil {
		f.logger.Error("error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames, partStatuses string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &fwp.IsActive, &partIDs, &partNames, &partStatuses)
		if err != nil {
			f.logger.Errorf("error: %v", err.Error())
			return []model.FacilityWithParts{}, err
		}

		if partIDs != "" {
			ids := strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			statuses := strings.Split(partStatuses, ",")
			for i, idStr := range ids {
				var status bool
				id, _ := strconv.Atoi(idStr)
				switch statuses[i] {
				case "1":
					status = true
					break
				case "0":
					status = false
				}
				fwp.Parts = append(fwp.Parts, model.Part{ID: id, Name: names[i], IsActive: status})
			}
		}

		facilities = append(facilities, fwp)
	}
	return facilities, nil
}

func (f FacilityRepository) GetByID(ctx context.Context, facilityID int) ([]model.FacilityWithParts, error) {

	baseQuery := `SELECT facility.id, facility.name, facility.have_parts, facility.is_active,
                  COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
                  COALESCE(GROUP_CONCAT(part.name), '') AS part_names,
                  COALESCE(GROUP_CONCAT(part.is_active), '') AS part_is_active
           FROM facility
           LEFT JOIN part ON facility.id = part.facility_id
           LEFT JOIN application ON application.facility_id = facility.id WHERE facility.id = $1 GROUP BY facility.id`

	rows, err := f.db.QueryContext(ctx, baseQuery, facilityID)
	if err != nil {
		f.logger.Error("error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames, partStatuses string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &fwp.IsActive, &partIDs, &partNames, &partStatuses)
		if err != nil {
			f.logger.Errorf("error: %v", err.Error())
			return []model.FacilityWithParts{}, err
		}

		if partIDs != "" {
			ids := strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			statuses := strings.Split(partStatuses, ",")
			for i, idStr := range ids {
				var status bool
				id, _ := strconv.Atoi(idStr)
				switch statuses[i] {
				case "1":
					status = true
					break
				case "0":
					status = false
				}
				fwp.Parts = append(fwp.Parts, model.Part{ID: id, Name: names[i], IsActive: status})
			}
		}

		facilities = append(facilities, fwp)
	}
	return facilities, nil
}
func (f FacilityRepository) GetActive(ctx context.Context, categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error) {
	args := make([]interface{}, 0, 3)

	baseQuery := `SELECT facility.id, facility.name, facility.have_parts, facility.is_active,
                  COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
                  COALESCE(GROUP_CONCAT(part.name), '') AS part_names,
                  COALESCE(GROUP_CONCAT(part.is_active), '') AS part_is_active
           FROM facility
           LEFT JOIN part ON facility.id = part.facility_id
           LEFT JOIN application ON application.facility_id = facility.id
           LEFT JOIN work_type ON application.work_type_id = work_type.id
           LEFT JOIN events ON application.event_id = events.id
           LEFT JOIN details ON events.details_id = details.id
           WHERE facility.is_active = TRUE`

	if categoryName != "" {
		baseQuery += ` AND details.category = ?`
		args = append(args, categoryName)
	}
	if workTypeID != 0 {
		baseQuery += ` AND work_type_id.id = ?`
		args = append(args, workTypeID)
	}
	if status != "" {
		baseQuery += ` AND application.status = ?`
		args = append(args, status)
	}

	baseQuery += " GROUP BY facility.id;"

	rows, err := f.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		f.logger.Error("error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames, partStatuses string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &fwp.IsActive, &partIDs, &partNames, &partStatuses)
		if err != nil {
			f.logger.Errorf("error: %v", err.Error())
			return []model.FacilityWithParts{}, err
		}

		if partIDs != "" {
			ids := strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			statuses := strings.Split(partStatuses, ",")
			for i, idStr := range ids {
				var status bool
				id, _ := strconv.Atoi(idStr)
				switch statuses[i] {
				case "1":
					status = true
					break
				case "0":
					status = false
				}
				fwp.Parts = append(fwp.Parts, model.Part{ID: id, Name: names[i], IsActive: status})
			}
		}

		facilities = append(facilities, fwp)
	}
	return facilities, nil
}
func (f FacilityRepository) GetByDate(ctx context.Context, startDate string, startTime string, endDate string, endTime string, facilityID int, bookingID int) ([]model.FacilityWithParts, error) {
	query := `
    SELECT 
		facility.id, 
		facility.name, 
		facility.have_parts,
		COALESCE(GROUP_CONCAT(DISTINCT CASE WHEN b.id IS NULL THEN part.id END), '') AS part_ids,
		COALESCE(GROUP_CONCAT(DISTINCT CASE WHEN b.id IS NULL THEN part.name END), '') AS part_names
	FROM facility
	LEFT JOIN part ON facility.id = part.facility_id AND part.is_active = TRUE
	LEFT JOIN booking_part bp ON part.id = bp.part_id
	LEFT JOIN booking b ON bp.booking_id = b.id 
		AND (
			(datetime(b.start_date || ' ' || b.start_time) <= datetime($1 || ' ' || $2) AND datetime(b.end_date || ' ' || b.end_time) >= datetime($1 || ' ' || $2)) 
			OR (datetime(b.start_date || ' ' || b.start_time) <= datetime($3 || ' ' || $4) AND datetime(b.end_date || ' ' || b.end_time) >= datetime($3 || ' ' || $4)) 
			OR (datetime(b.start_date || ' ' || b.start_time) >= datetime($1 || ' ' || $2) AND datetime(b.end_date || ' ' || b.end_time) <= datetime($3 || ' ' || $4))
		)
	WHERE facility.is_active = TRUE
	AND facility.id <> $5 -- Exclude the facility with the given facilityID
	GROUP BY facility.id
	HAVING (facility.have_parts = FALSE AND COUNT(DISTINCT b.id) = 0) 
		OR (facility.have_parts = TRUE AND COUNT(DISTINCT part.id) > COUNT(DISTINCT CASE WHEN b.id IS NOT NULL THEN part.id END))
    `

	rows, err := f.db.QueryContext(ctx, query, startDate, startTime, endDate, endTime, facilityID)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &partIDs, &partNames)
		if err != nil {
			f.logger.Errorf("error: %v", err.Error())
			return nil, err
		}

		if partIDs != "" {
			ids := strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			for i, idStr := range ids {
				id, _ := strconv.Atoi(idStr)
				fwp.Parts = append(fwp.Parts, model.Part{ID: id, Name: names[i]})
			}
		}

		facilities = append(facilities, fwp)
	}

	if facilityID != 0 && bookingID != 0 {
		additionalQuery := `SELECT 
			facility.id, 
			facility.name, 
			facility.have_parts,
			COALESCE(GROUP_CONCAT(DISTINCT CASE 
				WHEN b.id IS NULL OR b.id = $5 THEN part.id 
				END), '') AS part_ids,
			COALESCE(GROUP_CONCAT(DISTINCT CASE 
				WHEN b.id IS NULL OR b.id = $5 THEN part.name 
				END), '') AS part_names
		FROM facility
		LEFT JOIN part ON facility.id = part.facility_id AND part.is_active = TRUE
		LEFT JOIN booking_part bp ON part.id = bp.part_id
		LEFT JOIN booking b ON bp.booking_id = b.id 
			AND (
				(datetime(b.start_date || ' ' || b.start_time) <= datetime($1 || ' ' || $2) AND datetime(b.end_date || ' ' || b.end_time) >= datetime($1 || ' ' || $2)) 
				OR (datetime(b.start_date || ' ' || b.start_time) <= datetime($3 || ' ' || $4) AND datetime(b.end_date || ' ' || b.end_time) >= datetime($3 || ' ' || $4)) 
				OR (datetime(b.start_date || ' ' || b.start_time) >= datetime($1 || ' ' || $2) AND datetime(b.end_date || ' ' || b.end_time) <= datetime($3 || ' ' || $4))
			)
		WHERE facility.is_active = TRUE AND facility.id = $6
		GROUP BY facility.id
		HAVING (facility.have_parts = FALSE AND COUNT(DISTINCT b.id) = 0) 
			OR (facility.have_parts = TRUE AND COUNT(DISTINCT part.id) > COUNT(DISTINCT CASE WHEN b.id IS NOT NULL THEN part.id END))`

		additionalRows, err := f.db.QueryContext(ctx, additionalQuery, startDate, startTime, endDate, endTime, bookingID, facilityID)
		if err != nil {
			f.logger.Error("error: ", err.Error())
			return nil, err
		}
		defer additionalRows.Close()

		for additionalRows.Next() {
			var fwp model.FacilityWithParts
			var partIDs, partNames string

			err = additionalRows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &partIDs, &partNames)
			if err != nil {
				f.logger.Error("error: ", err.Error())
				return nil, err
			}

			if partIDs != "" {
				ids := strings.Split(partIDs, ",")
				names := strings.Split(partNames, ",")
				for i, idStr := range ids {
					id, _ := strconv.Atoi(idStr)
					fwp.Parts = append(fwp.Parts, model.Part{ID: id, Name: names[i]})
				}
			}

			facilities = append(facilities, fwp)
		}
	}

	return facilities, nil
}

// time format should be HH:MM:SS
func (f FacilityRepository) GetByDateTime(ctx context.Context, startDate string, startTime string, endDate string, endTime string) ([]model.FacilityWithParts, error) {
	//startTime += ":00"
	//endTime += ":00"

	ctx = context.Background()

	var facilities []model.FacilityWithParts

	facilityRows, err := f.db.QueryContext(ctx, `SELECT facility.id, facility.name, facility.have_parts, facility.is_active
												FROM facility
												WHERE NOT EXISTS(SELECT 1
												                 FROM booking
												                 LEFT JOIN booking_part  on booking.id = booking_part.booking_id
												                 WHERE 
												                     booking.facility_id = facility.id
												                	AND 
												                     (
												                     	(booking.start_date || ' ' || booking.start_time) BETWEEN $1 || ' ' || $2 AND $3 || ' ' || $4
												                 	OR	(booking.end_date || ' ' || booking.end_time) BETWEEN $1 || ' ' || $2 AND $3 || ' ' || $4
												                     )
												                 GROUP BY booking.id
												                 HAVING COUNT(booking_part.booking_id) = 0);`, startDate, startTime, endDate, endTime)

	if err != nil {
		f.logger.Errorf(err.Error())
		return []model.FacilityWithParts{}, err
	}

	defer facilityRows.Close()

	for facilityRows.Next() {
		var facility model.FacilityWithParts

		err = facilityRows.Scan(&facility.ID,
			&facility.Name,
			&facility.HaveParts,
			&facility.IsActive,
		)

		if facility.HaveParts {
			var parts []model.Part

			partRows, err := f.db.QueryContext(ctx, `SELECT part.id, part.name, part.facility_id, part.is_active 
													FROM part 
													WHERE 
													    part.facility_id = $1 
													AND
													    NOT EXISTS(SELECT 1 
												                 FROM booking_part 
												                 INNER JOIN booking ON booking_part.booking_id = booking.id
												                 WHERE booking_part.part_id = part.id  
												                	AND 
												                     (
												                     	(booking.start_date || ' ' || booking.start_time) BETWEEN $2 || ' ' || $3 AND $4 || ' ' || $5
												                 	OR	(booking.end_date || ' ' || booking.end_time) BETWEEN $2 || ' ' || $3 AND $4 || ' ' || $5
												                     )
												                     );`, facility.ID, startDate, startTime, endDate, endTime)
			if err != nil {
				f.logger.Errorf(err.Error())
				return []model.FacilityWithParts{}, err
			}

			defer partRows.Close()

			for partRows.Next() {
				var part model.Part

				err = partRows.Scan(&part.ID,
					&part.Name,
					&part.FacilityID,
					&part.IsActive,
				)

				if err != nil {
					f.logger.Errorf(err.Error())
					return []model.FacilityWithParts{}, err
				}
				parts = append(parts, part)
			}
			if len(parts) != 0 {
				facility.Parts = parts
				facilities = append(facilities, facility)
			} else {
				// не добавляем facility, у которого have parts True, но все part заняты
				_ = facility
			}
		} else {
			// если у facility have parts False, то добавляем естественно I love hot dogs
			facilities = append(facilities, facility)
		}
	}

	return facilities, nil
}

func (f FacilityRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility name update", int(rowsCount))
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Error("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (f FacilityRepository) Delete(ctx context.Context, facilityId int, isActive bool) error {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		f.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET is_active = $1 WHERE facility.id = $2;`, isActive, facilityId)
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility delete", int(rowsCount))
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
