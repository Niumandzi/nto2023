package model

type Details struct {
	ID       int
	Category string
	TypeName string
	IsActive bool
}

type Event struct {
	ID          int
	Name        string
	Date        string
	Description string
	IsActive    bool
	DetailsID   int
}

type EventWithDetails struct {
	ID          int
	Name        string
	Date        string
	Description string
	IsActive    bool
	Details     Details
}
