package work

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s WorkTypeService) UpdateWorkType(workTypeID int, name string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.workTypeRepo.Update(ctx, workTypeID, name)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
