package model

type Category struct {
	ID           int
	CategoryType string
}

type EventType struct {
	ID         int
	CategoryID int
	EventType  string
}

type Event struct {
	ID          int
	EventTypeId int
	Date        string
	Name        string
	Description string
}

type EventWithCategoryAndType struct {
	ID          int
	EventType   EventType
	Date        string
	Name        string
	Description string
}
