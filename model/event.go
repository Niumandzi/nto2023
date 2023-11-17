package model

type Details struct {
	Id       int
	TypeName string
	Category string
}

type Event struct {
	Id          int
	Name        string
	Description string
	Date        string
	DetailsId   int
}

type EventWithDetails struct {
	Id          int
	Name        string
	Description string
	Date        string
	Details     Details
}
