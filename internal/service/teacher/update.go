package teacher

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (t TeacherService) UpdateTeacher(teacherID int, name string) error {
	ctx, cancel := context.WithTimeout(t.ctx, t.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		t.logger.Error("error: %v", err)
		return err
	}

	err = t.teacherRepo.Update(ctx, teacherID, name)
	if err != nil {
		t.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
