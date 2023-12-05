package model

type Booking struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	StartTime   string
	EndDate     string
	EndTime     string
	EventID     int
	FacilityID  int
	PartIDs     []int
}

type BookingWithFacility struct {
	ID          int
	Description string
	CreateDate  string
	StartDate   string
	StartTime   string
	EndDate     string
	EndTime     string
	Event       EventWithDetails
	Facility    Facility
	Parts       []Part
}
