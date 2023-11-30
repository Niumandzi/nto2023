package event

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s DetailsService) UpdateDetail(detailsID int, typeName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.Validate(typeName, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return err
	}

	err = s.detailsRepo.UpdateTypeName(ctx, detailsID, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
