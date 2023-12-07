package model

type MugType struct {
	ID       int
	Name     string
	IsActive bool
}

type Teacher struct {
	ID       int
	Name     string
	IsActive bool
}

type Schedule struct {
	ID        int
	Name      int
	StartTime string
	EndTime   string
}

type Registration struct {
	ID           int
	Name         string
	StartDate    string
	NumberOfDays int
	FacilityID   int
	MugTypeID    int
	TeacherID    int
	Schedule     []Schedule
	PartIDs      []int
}

type RegistrationWithDetails struct {
	ID           int
	Name         string
	StartDate    string
	NumberOfDays string
	Facility     Facility
	MugType      MugType
	Teacher      Teacher
	Schedule     []Schedule
	Parts        []Part
}
