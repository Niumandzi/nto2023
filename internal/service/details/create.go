package event

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/net/context"
)

func (s DetailsService) CreateDetail(categoryName string, typeName string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	err := validation.Validate(typeName, validation.Required)
	if err != nil {
		s.logger.Error("error: %v", err)
		return 0, err
	}

	id, err := s.detailsRepo.Create(ctx, categoryName, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, err
}
