package facility

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
		f.logger.Error("error: ", err.Error())
		return 0, err
	}

	haveParts := len(parts) > 0
	res, err := tx.ExecContext(ctx, `INSERT INTO facility (name, have_parts) VALUES ($1, $2);`, name, haveParts)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	if haveParts {
		for _, partName := range parts {
			_, err := tx.ExecContext(ctx, `INSERT INTO part (name, facility_id) VALUES ($1, $2);`, partName, id)
			if err != nil {
				f.logger.Error("error: ", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (f FacilityRepository) Get(ctx context.Context, categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error) {
	args := make([]interface{}, 0, 2)
	kwargs := make(map[string]interface{})
	var query string

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
           WHERE `

	if categoryName != "" {
		kwargs["details.category"] = categoryName
	}
	if workTypeID != 0 {
		kwargs["work_type.id"] = workTypeID
	}
	if status != "" {
		kwargs["application.status"] = status
	}
	if categoryName == "" && workTypeID == 0 && status == "" {
		baseQuery = `SELECT facility.id, facility.name, facility.have_parts, facility.is_active,
          			COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
       				COALESCE(GROUP_CONCAT(part.name), '') AS part_names,
                    COALESCE(GROUP_CONCAT(part.is_active), '') AS part_is_active
		  FROM facility
          LEFT JOIN part ON facility.id = part.facility_id
          GROUP BY facility.id;`
	}

	length := len(kwargs)

	if length != 0 {
		i := 0
		for key, val := range kwargs {
			if i == length-1 {
				query += fmt.Sprintf("%v = ?;", key)
				args = append(args, val)
			} else {
				query += fmt.Sprintf("%v = ? AND ", key)
				args = append(args, val)
			}
			i++
		}
		query += " GROUP BY facility.id;"
	}

	rows, err := f.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames, partStatuses string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &fwp.IsActive, &partIDs, &partNames, &partStatuses)
		if err != nil {
			f.logger.Error("error: ", err.Error())
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

func (f FacilityRepository) GetByDate(ctx context.Context, startDate string, endDate string) ([]model.FacilityWithParts, error) {
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
	LEFT JOIN booking b ON (bp.booking_id = b.id AND ((b.start_date <= $1 AND b.end_date >= $1) OR (b.start_date <= $2 AND b.end_date >= $2) OR (b.start_date >= $1 AND b.end_date <= $2)))
	WHERE facility.is_active = TRUE
	GROUP BY facility.id
	HAVING (facility.have_parts = FALSE AND COUNT(DISTINCT b.id) = 0) OR (facility.have_parts = TRUE AND COUNT(DISTINCT CASE WHEN b.id IS NOT NULL THEN b.id END) < COUNT(DISTINCT part.id))
    `

	rows, err := f.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &partIDs, &partNames)
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

	return facilities, nil
}

func (f FacilityRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility name update", int(rowsCount))
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (f FacilityRepository) Delete(ctx context.Context, facilityId int, isActive bool) error {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET is_active = $1 WHERE facility.id = $2;`, isActive, facilityId)
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility delete", int(rowsCount))
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		f.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
