package mug

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (m MugTypeService) UpdateMugType(mugTypeID int, name string) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		m.logger.Error("error: %v", err)
		return err
	}

	err = m.mugTypeRepo.Update(ctx, mugTypeID, name)
	if err != nil {
		m.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
