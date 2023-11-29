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

func (w FacilityRepository) Create(ctx context.Context, name string, parts []string) (int, error) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return 0, err
	}

	haveParts := len(parts) > 0 // Устанавливаем haveParts в true, если массив parts не пустой

	res, err := tx.ExecContext(ctx, `INSERT INTO facility (name, have_parts) VALUES ($1, $2);`, name, haveParts)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	if haveParts {
		for _, partName := range parts {
			_, err := tx.ExecContext(ctx, `INSERT INTO part (name, facility_id) VALUES ($1, $2);`, partName, id)
			if err != nil {
				w.logger.Errorf("error: %v", err.Error())
				tx.Rollback()
				return 0, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (w FacilityRepository) Get(ctx context.Context, categoryName string, workTypeID int, status string) ([]model.FacilityWithParts, error) {
	args := make([]interface{}, 0, 2)
	kwargs := make(map[string]interface{})
	var query string

	//query = `SELECT DISTINCT facility.id, facility.name, facility.have_parts
	//      FROM facility
	//      INNER JOIN application ON application.facility_id = facility.id
	//      INNER JOIN work_type ON application.work_type_id = work_type.id
	//      INNER JOIN events ON application.event_id = events.id
	//      INNER JOIN details ON events.details_id = details.id
	//      WHERE `

	query = `SELECT facility.id, facility.name, facility.have_parts,
	     			COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
       				COALESCE(GROUP_CONCAT(part.name), '') AS part_names
		 FROM facility
	     LEFT JOIN part ON facility.id = part.facility_id
	     INNER JOIN application ON application.facility_id = facility.id
	     INNER JOIN work_type ON application.work_type_id = work_type.id
	     INNER JOIN events ON application.event_id = events.id
	     INNER JOIN details ON events.details_id = details.id
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
		query = `SELECT facility.id, facility.name, facility.have_parts, 
          			COALESCE(GROUP_CONCAT(part.id), '') AS part_ids,
       				COALESCE(GROUP_CONCAT(part.name), '') AS part_names
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

	rows, err := w.db.QueryContext(ctx, query, args...)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return nil, err
	}

	defer rows.Close()

	var facilities []model.FacilityWithParts
	for rows.Next() {
		var fwp model.FacilityWithParts
		var partIDs, partNames string

		err = rows.Scan(&fwp.ID, &fwp.Name, &fwp.HaveParts, &partIDs, &partNames)
		if err != nil {
			w.logger.Errorf("error: %v", err.Error())
			return []model.FacilityWithParts{}, err
		}

		if partIDs != "" {
			ids := strings.Split(partIDs, ",")
			names := strings.Split(partNames, ",")
			for i, idStr := range ids {
				id, _ := strconv.Atoi(idStr)
				fwp.Parts = append(fwp.Parts, model.Parts{ID: id, Name: names[i]})
			}
		}

		facilities = append(facilities, fwp)
	}

	return facilities, nil
}

func (w FacilityRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE facility SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility name update", int(rowsCount))
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (w FacilityRepository) Delete(ctx context.Context, id int) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		w.logger.Error("error: %v", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM facility WHERE facility.id = $1;`, id)
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errors.NewRowCountError("facility delete", int(rowsCount))
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		w.logger.Errorf("error: %v", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
