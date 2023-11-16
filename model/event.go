package model

type EventType struct {
	ID        int
	EventType string
}

type Category struct {
	ID          int
	EventTypeID int
	Category    string
}

type CategoryWithEventType struct {
	ID        int
	Category  string
	EventType EventType
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
	Date        string
	Name        string
	Description string
	Category    CategoryWithEventType
}
