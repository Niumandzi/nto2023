package mug

import (
	"context"
	"database/sql"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type MugTypeRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NemMugTypeRepository(db *sql.DB, logger logging.Logger) MugTypeRepository {
	return MugTypeRepository{
		db:     db,
		logger: logger,
	}
}

func (m MugTypeRepository) Create(ctx context.Context, name string) (int, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO mug_type (name) VALUES ($1);`, name)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (m MugTypeRepository) Get(ctx context.Context) ([]model.MugType, error) {
	var mugTypes []model.MugType

	baseQuery := `SELECT mug_type.id, mug_type.name, mug_type.is_active FROM mug_type`

	roms, err := m.db.QueryContext(ctx, baseQuery)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer roms.Close()

	for roms.Next() {
		var mugType model.MugType

		err = roms.Scan(&mugType.ID, &mugType.Name, &mugType.IsActive)
		if err != nil {
			m.logger.Error("error: ", err.Error())
			return nil, err
		}

		mugTypes = append(mugTypes, mugType)
	}

	return mugTypes, nil
}

func (m MugTypeRepository) GetActive(ctx context.Context, facilityID int, teacherID int) ([]model.MugType, error) {
	args := make([]interface{}, 0, 2)
	var mugTypes []model.MugType

	baseQuery := `SELECT mug_type.id,
                         mug_type.name
			FROM mug_type
			LEFT JOIN registration ON registration.mug_type_id = mug_type.id
			LEFT JOIN facility ON registration.facility_id = facility.id
			LEFT JOIN teacher ON registration.teacher_id = teacher.id
			WHERE mug_type.is_active = TRUE`

	if facilityID != 0 {
		baseQuery += ` AND facility.id = ?`
		args = append(args, facilityID)
	}
	if teacherID != 0 {
		baseQuery += ` AND teacher.id = ?`
		args = append(args, teacherID)
	}

	roms, err := m.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer roms.Close()

	for roms.Next() {
		var mugType model.MugType

		err = roms.Scan(&mugType.ID, &mugType.Name)
		if err != nil {
			m.logger.Error("error: ", err.Error())
			return nil, err
		}

		mugTypes = append(mugTypes, mugType)
	}

	return mugTypes, nil
}

func (m MugTypeRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE mug_type SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	romsCount, err := res.RowsAffected()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if romsCount != 1 {
		err = errcode.NewRowCountError("mug type name update", int(romsCount))
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (m MugTypeRepository) Delete(ctx context.Context, mugTypeID int, isActive bool) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE mug_type SET is_active = $1 WHERE mug_type.id = $2;`, isActive, mugTypeID)
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	romsCount, err := res.RowsAffected()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if romsCount != 1 {
		err = errcode.NewRowCountError("mug type delete", int(romsCount))
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		m.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
