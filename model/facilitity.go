package model

type Part struct {
	ID         int
	FacilityID int
	Name       string
}

type Facility struct {
	ID         int
	Name       string
	IsTwoParts bool
}
