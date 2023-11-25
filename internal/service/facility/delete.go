package facility

import "context"

func (s FacilityService) DeleteFacility(id int) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)

	defer cancel()

	err := s.facilityRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("error: %v", err.Error())
		return err
	}
	return nil
}
