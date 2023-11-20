package event

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (s DetailsService) UpdateDetail(detailsId int, typeName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := validation.ValidateStruct(&typeName,
		validation.Field(&typeName, validation.Required),
	)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	err = s.detailsRepo.UpdateTypeName(ctx, detailsId, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
