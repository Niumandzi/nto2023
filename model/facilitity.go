package model

type Parts struct {
	ID         int
	FacilityID int
	Name       string
}

type Facility struct {
	ID        int
	Name      string
	HaveParts bool
}

type FacilityWithParts struct {
	ID        int
	Name      string
	HaveParts bool
	Parts     []Parts
}
