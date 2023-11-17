package errors

import "fmt"

type RowCountError struct {
	RowsCount     int
	WhereHappened string
}

func (r RowCountError) Error() string {
	return fmt.Sprintf("where happened: %v, how many rows: %v", r.WhereHappened, r.RowsCount)
}

func NewRowCountError(where string, count int) *RowCountError {
	return &RowCountError{
		RowsCount:     count,
		WhereHappened: where,
	}
}
