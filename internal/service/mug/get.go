package mug

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (m MugTypeService) GetMugTypes() ([]model.MugType, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.contextTimeout)
	defer cancel()

	types, err := m.mugTypeRepo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return types, nil
}

func (m MugTypeService) GetActiveMugTypes(facilityID int, teacherID int) ([]model.MugType, error) {
	ctx, cancel := context.WithTimeout(m.ctx, m.contextTimeout)
	defer cancel()

	types, err := m.mugTypeRepo.GetActive(ctx, facilityID, teacherID)
	if err != nil {
		return nil, err
	}

	return types, nil
}
