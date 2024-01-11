package mug

import "context"

func (m MugTypeService) DeleteRestoreMugType(mugTypeID int, isActive bool) error {
	ctx, cancel := context.WithTimeout(m.ctx, m.contextTimeout)
	defer cancel()

	err := m.mugTypeRepo.Delete(ctx, mugTypeID, isActive)
	if err != nil {
		m.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
