package work

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/niumandzi/nto2023/model"
)

func (s WorkTypeService) UpdateWorkType(workType model.WorkType) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.Validate(workType.Name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.workTypeRepo.Update(ctx, workType.ID, workType.Name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
