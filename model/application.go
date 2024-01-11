package model

type WorkType struct {
	ID       int
	Name     string
	IsActive bool
}

type Application struct {
	ID          int
	Description string
	CreateDate  string
	DueDate     string
	Status      string
	WorkTypeId  int
	FacilityId  int
	EventId     int
}

type ApplicationWithDetails struct {
	ID          int
	WorkType    WorkType
	Facility    Facility
	Description string
	CreateDate  string
	DueDate     string
	Status      string
	Event       EventWithDetails
}
