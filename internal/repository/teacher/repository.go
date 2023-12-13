package teacher

import (
	"context"
	"database/sql"
	errcode "github.com/niumandzi/nto2023/internal/errors"
	"github.com/niumandzi/nto2023/model"
	"github.com/niumandzi/nto2023/pkg/logging"
)

type TeacherRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewTeacherRepository(db *sql.DB, logger logging.Logger) TeacherRepository {
	return TeacherRepository{
		db:     db,
		logger: logger,
	}
}

func (t TeacherRepository) Create(ctx context.Context, name string) (int, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `INSERT INTO teacher (name) VALUES ($1);`, name)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}

func (t TeacherRepository) Get(ctx context.Context) ([]model.Teacher, error) {
	var teachers []model.Teacher

	baseQuery := `SELECT teacher.id, teacher.name, teacher.is_active FROM teacher`

	rows, err := t.db.QueryContext(ctx, baseQuery)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var teacher model.Teacher

		err = rows.Scan(&teacher.ID, &teacher.Name, &teacher.IsActive)
		if err != nil {
			t.logger.Error("error: ", err.Error())
			return nil, err
		}

		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (t TeacherRepository) GetActive(ctx context.Context, facilityID int, mugTypeID int) ([]model.Teacher, error) {
	args := make([]interface{}, 0, 2)
	var teachers []model.Teacher

	baseQuery := `SELECT teacher.id,
                         teacher.name
			FROM teacher
			LEFT JOIN registration ON registration.teacher_id = teacher.id
			LEFT JOIN facility ON registration.facility_id = facility.id
			LEFT JOIN mug_type ON registration.mug_type_id = mug_type.id
			WHERE teacher.is_active = TRUE`

	if facilityID != 0 {
		baseQuery += ` AND facility.id = ?`
		args = append(args, facilityID)
	}
	if mugTypeID != 0 {
		baseQuery += ` AND mug_type.id = ?`
		args = append(args, mugTypeID)
	}

	rows, err := t.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var teacher model.Teacher

		err = rows.Scan(&teacher.ID, &teacher.Name)
		if err != nil {
			t.logger.Error("error: ", err.Error())
			return nil, err
		}

		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func (t TeacherRepository) Update(ctx context.Context, idOld int, nameUpd string) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE teacher SET name = $1 WHERE id = $2;`, nameUpd, idOld)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("teacher name update", int(rowsCount))
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (t TeacherRepository) Delete(ctx context.Context, teacherID int, isActive bool) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM teacher WHERE teacher.id = $1;`, teacherID)
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}
	if rowsCount != 1 {
		err = errcode.NewRowCountError("teacher delete", int(rowsCount))
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		t.logger.Error("error: ", err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
