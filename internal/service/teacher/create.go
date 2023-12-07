package teacher

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/net/context"
)

func (t TeacherService) CreateTeacher(name string) (int, error) {
	ctx, cancel := context.WithTimeout(t.ctx, t.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		t.logger.Error("error: %v", err)
		return 0, err
	}

	id, err := t.teacherRepo.Create(ctx, name)
	if err != nil {
		t.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, err
}
