package part

import "context"

func (s PartService) CreatePart(facilityID int, partNames []string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	id, err := s.partRepo.Create(ctx, facilityID, partNames)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return 0, err
	}

	return id, nil
}
