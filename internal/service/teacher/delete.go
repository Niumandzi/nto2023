package teacher

import "context"

func (t TeacherService) DeleteRestoreTeacher(teacherID int, isActive bool) error {
	ctx, cancel := context.WithTimeout(t.ctx, t.contextTimeout)
	defer cancel()

	err := t.teacherRepo.Delete(ctx, teacherID, isActive)
	if err != nil {
		t.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
