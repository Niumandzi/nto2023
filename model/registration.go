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
	ID           int
	NumberOfDays int
	Monday       bool
	Tuesday      bool
	Wednesday    bool
	Thursday     bool
	Friday       bool
	Saturday     bool
	Sunday       bool
}

type Registration struct {
	ID         int
	Name       WorkType
	StartDate  Facility
	ScheduleID string
	FacilityID string
	MugTypeID  string
	TeacherID  string
}

type RegistrationWithDetails struct {
	ID        int
	Name      WorkType
	StartDate Facility
	Schedule  Schedule
	Facility  FacilityWithParts
	MugType   MugType
	Teacher   Teacher
}
