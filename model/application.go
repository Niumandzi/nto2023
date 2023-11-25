package model

type WorkType struct {
	ID   int
	Name string
}

type Facility struct {
	ID   int
	Name string
}

type Application struct {
	ID          int
	Description string
	CreateDate  string
	DueDate     string
	Status      string
	WorkTypeId  int
	EventId     int
	FacilityId  int
}

type ApplicationWithDetails struct {
	ID          int
	Description string
	CreateDate  string
	DueDate     string
	Status      string
	WorkType    WorkType
	Event       EventWithDetails
	Facility    Facility
}
