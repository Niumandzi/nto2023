package model

type Facility struct {
	ID         int
	Name       string
	IsTwoParts bool
}

type Part struct {
	ID         int
	FacilityID int
	Name       string
}
