package model

type Booking struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	EndDate     string
	EventID     string
	FacilityID  int
	PartID      int
}

type BookingWithFacility struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	EndDate     string
	EventID     string
	Facility    Facility
	Part        Parts
}
