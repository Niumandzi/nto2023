package event

import (
	"golang.org/x/net/context"
)

func (s DetailsService) CreateDetail(categoryName string, typeName string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	id, err := s.detailsRepo.Create(ctx, categoryName, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, err
}
