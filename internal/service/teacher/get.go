package teacher

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (t TeacherService) GetTeachers() ([]model.Teacher, error) {
	ctx, cancel := context.WithTimeout(t.ctx, t.contextTimeout)
	defer cancel()

	types, err := t.teacherRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (t TeacherService) GetActiveTeachers(facilityID int, mugTypeID int) ([]model.Teacher, error) {
	ctx, cancel := context.WithTimeout(t.ctx, t.contextTimeout)
	defer cancel()

	teachers, err := t.teacherRepo.GetActive(ctx, facilityID, mugTypeID)
	if err != nil {
		return nil, err
	}

	return teachers, nil
}
