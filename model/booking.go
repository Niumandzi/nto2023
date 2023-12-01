package model

type Booking struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	EndDate     string
	EventID     int
	FacilityID  int
	PartIDs     []int
}

type BookingWithFacility struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	EndDate     string
	EventID     int
	Facility    Facility
	Parts       []Parts
}
