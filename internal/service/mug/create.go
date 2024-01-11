package mug

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/net/context"
)

func (m MugTypeService) CreateMugType(name string) (int, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.contextTimeout)
	defer cancel()

	err := validation.Validate(name, validation.Required)
	if err != nil {
		m.logger.Error("error: %v", err)
		return 0, err
	}

	id, err := m.mugTypeRepo.Create(ctx, name)
	if err != nil {
		m.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, err
}
