package validations

import (
	"errors"
	"regexp"
	"strconv"
)

var basicDateRegex = regexp.MustCompile(`^(\d{2})\.(\d{2})\.(\d{4})$`)

func ValidateDate(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for date")
	}

	matches := basicDateRegex.FindStringSubmatch(s)
	if matches == nil {
		return errors.New("invalid date format, expected dd.mm.yyyy")
	}

	day, _ := strconv.Atoi(matches[1])
	month, _ := strconv.Atoi(matches[2])
	year, _ := strconv.Atoi(matches[3])

	if month < 1 || month > 12 {
		return errors.New("month must be between 01 and 12")
	}
	if year < 2000 || year > 2099 {
		return errors.New("year must be between 2000 and 2099")
	}

	if !isValidDay(day, month, year) {
		return errors.New("invalid number of days for the given month and year")
	}

	return nil
}

func isValidDay(day, month, year int) bool {
	daysInMonth := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if month == 2 && (year%400 == 0 || (year%100 != 0 && year%4 == 0)) {
		daysInMonth[1] = 29
	}

	return day >= 1 && day <= daysInMonth[month-1]
}
