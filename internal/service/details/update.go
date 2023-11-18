package event

import (
	"context"
)

func (s DetailsService) UpdateDetail(detailsId int, typeName string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.UpdateTypeName(ctx, detailsId, typeName)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}

	return nil
}
