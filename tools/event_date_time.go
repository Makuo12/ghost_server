package tools

import (
	"fmt"
	"time"
)

func DatesMatch(startDate string, endDate string, dbStartDate time.Time, dbEndDate time.Time) (match bool) {
	if startDate == ConvertDateOnlyToString(dbStartDate) && endDate == ConvertDateOnlyToString(dbEndDate) {
		match = true

	}
	return
}

func DatesMatchString(startDate string, endDate string, newStartDate string, newEndDate string) (match bool) {
	if startDate == newStartDate && endDate == newEndDate {
		match = true
	}
	return
}

func RemoveRecurDate(date string, dates []string) (newDates []string, err error) {

	var dateFound bool
	for _, d := range dates {
		if d == date {
			dateFound = true
		} else {
			newDates = append(newDates, d)
		}
	}
	if !dateFound {
		err = fmt.Errorf("this event date cannot be found")
	}
	if len(newDates) == 0 {
		newDates = []string{"none"}
	}
	return
}

func UpdateRecurDate(newDate string, date string, dates []string) (newDates []string, err error) {

	var dateFound bool
	for _, d := range dates {
		if d == date {
			dateFound = true
		} else {
			newDates = append(newDates, d)
		}
	}
	newDates = append(newDates, newDate)
	if !dateFound {
		err = fmt.Errorf("this event date cannot be found")
	}
	return
}
