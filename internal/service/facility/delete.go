package facility

import "context"

func (s FacilityService) DeleteRestoreFacility(id int, isActive bool) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.facilityRepo.Delete(ctx, id, isActive)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
