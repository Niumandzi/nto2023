package event

import "context"

func (s DetailsService) DeleteDetail(detailsId int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.detailsRepo.DeleteType(ctx, detailsId)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
