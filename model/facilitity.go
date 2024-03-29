package model

type Part struct {
	ID         int
	FacilityID int
	Name       string
	IsActive   bool
}

type Facility struct {
	ID        int
	Name      string
	HaveParts bool
	IsActive  bool
}

type FacilityWithParts struct {
	ID        int
	Name      string
	HaveParts bool
	IsActive  bool
	Parts     []Part
}
