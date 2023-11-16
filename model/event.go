package model

type EventType struct {
	ID        int
	EventType string
}

type EventTypeWithCategory struct {
	ID           int
	EventTypeID  int
	CategoryType string
}

type Event struct {
	ID          int
	EventTypeId int
	Date        string
	Name        string
	Description string
}

type EventWithCategoryAndType struct {
	ID                    int
	Date                  string
	Name                  string
	Description           string
	EventTypeWithCategory EventTypeWithCategory
}
