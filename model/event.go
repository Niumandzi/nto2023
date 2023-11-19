package model

type Details struct {
	ID       int
	Category string
	TypeName string
}

type Event struct {
	ID          int
	Name        string
	Date        string
	Description string
	DetailsID   int
}

type EventWithDetails struct {
	ID          int
	Name        string
	Description string
	Date        string
	Details     Details
}
