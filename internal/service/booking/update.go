package booking

import (
	"context"
	"github.com/niumandzi/nto2023/model"
)

func (s BookingService) UpdateBooking(bookingWithFacilityUpd model.BookingWithFacility) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.contextTimeout)
	defer cancel()

	partIDs := make([]int, 0, len(bookingWithFacilityUpd.Parts))

	for _, part := range bookingWithFacilityUpd.Parts {
		partID := part.ID
		partIDs = append(partIDs, partID)
	}

	bookingUpd := model.Booking{
		ID:          bookingWithFacilityUpd.ID,
		Description: bookingWithFacilityUpd.Description,
		CreateDate:  bookingWithFacilityUpd.CreateDate,
		StartDate:   bookingWithFacilityUpd.StartDate,
		EndDate:     bookingWithFacilityUpd.EndDate,
		EventID:     bookingWithFacilityUpd.Event.ID,
		FacilityID:  bookingWithFacilityUpd.Facility.ID,
		PartIDs:     partIDs,
	}

	err := s.bookingRepo.Update(ctx, bookingUpd)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	// получаем айди старых партов, которые нужно удалить
	oldPartIDs, err := s.bookingPartRepo.GetParts(ctx, bookingUpd.ID)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	for _, oldPartID := range oldPartIDs {
		err = s.bookingPartRepo.Delete(ctx, bookingUpd.ID, oldPartID)
		if err != nil {
			s.logger.Errorf("error on oldPartID: %v, error: %v", oldPartID, err.Error())
			return err
		}
	}

	for _, partID := range partIDs {
		err = s.bookingPartRepo.Create(ctx, bookingUpd.ID, partID)
		if err != nil {
			s.logger.Errorf("error on partID: %v, error: %v", partID, err.Error())
			return err
		}
	}

	return nil
}
