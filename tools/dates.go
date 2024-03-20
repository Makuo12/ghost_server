package tools

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func ConvertDateOnlyStringToDate(s string) (time.Time, error) {
	layout := "2006-01-02"
	t2, err := time.Parse(layout, s)
	return t2, err
}
func ConvertTimeToStringDateOnly(t time.Time) string {
	return t.Format("2006-01-02")
}

func CurrentDatePlusDaysToString(daysToAdd int) string {
	// Get the current date and time
	currentTime := time.Now()

	// Add the specified number of days to the current date
	newDate := currentTime.AddDate(0, 0, daysToAdd)

	// Format the new date as a string in the "2002-04-04" format
	dateString := newDate.Format("2006-01-02")

	return dateString
}

func CurrentDatePlusMonthsToString(monthsToAdd int) string {
	// Get the current date and time
	currentTime := time.Now()

	// Add the specified number of months to the current date (0 days added)
	newDate := currentTime.AddDate(0, monthsToAdd, 0)

	// Format the new date as a string in the "2002-04-04" format
	dateString := newDate.Format("2006-01-02")

	return dateString
}

func ConvertTimeToYear(t time.Time) string {
	date := ConvertTimeToStringDateOnly(t)
	return strings.Split(date, "-")[0]
}

// this convert a time to a string. If the time is a fake date it return an empty string
func ConvertDateOnlyToString(t time.Time) string {
	var s string = t.Format("2006-01-02")
	if s == FakeDate {
		return ""
	}
	return s
}

func ConvertTimeOnlyToString(t time.Time) string {
	var timeString string = t.Format("15:04")
	return timeString
}

func ConvertTimeToString(t time.Time) string {
	var timeString string = t.Format("2006-01-02T15:04:05")
	return timeString
}
func ConvertStringToTime(timeString string) (t time.Time, err error) {
	layout := "2006-01-02T15:04:05"
	t, err = time.Parse(layout, timeString)
	if err != nil {
		log.Println("Error parsing date:", err)
		return
	}
	return
}

func ConvertDateTimeStringToTime(dateString string, timeString string) (parsedDateTime time.Time, err error) {
	dateTimeLayout := "2006-01-02 15:04"
	combinedString := dateString + " " + timeString
	parsedDateTime, err = time.Parse(dateTimeLayout, combinedString)
	if err != nil {
		return
	}
	return
}

func ConvertStringToTimeOnly(s string) (time.Time, error) {
	layout := "15:04"
	t2, err := time.Parse(layout, s)
	return t2, err
}

func CheckItemIsNew(t time.Time) bool {
	t_1 := time.Now().Local().UTC()
	duration := t_1.Sub(t)
	// time.Hour*2400 is 100 days
	if duration < time.Duration(time.Hour*2400) {
		return true
	} else {
		return false
	}
}
func GetEventDateTime(startDate string, endDate string, startTime string, endTime string) (startingDate time.Time, endingDate time.Time, startingTime time.Time, endingTime time.Time, err error) {
	startingTime, err = ConvertStringToTimeOnly(startTime)
	if err != nil {
		return
	}
	endingTime, err = ConvertStringToTimeOnly(endTime)
	if err != nil {
		return
	}
	startingDate, err = ConvertDateOnlyStringToDate(startDate)
	if err != nil {
		return
	}
	endingDate, err = ConvertDateOnlyStringToDate(endDate)
	if err != nil {
		return
	}
	return
}

func GenerateDateListString(startDate, endDate string) ([]string, error) {
	if startDate == endDate {
		return []string{startDate}, nil
	}
	layout := "2006-01-02" // The format of the input dates
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	var dates []string
	current := start
	for current.Before(end) || current.Equal(end) {
		dates = append(dates, current.Format(layout))
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}

func GenerateDateListStringFromTime(startDate, endDate time.Time) []string {
	if startDate == endDate {
		return []string{ConvertDateOnlyToString(startDate)}
	}
	layout := "2006-01-02"
	if startDate.After(endDate) {
		return []string{}
	}

	var dates []string
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		dates = append(dates, current.Format(layout))
		current = current.AddDate(0, 0, 1)
	}

	return dates
}

func GenerateDateListTime(startDate, endDate string) ([]time.Time, error) {
	if startDate == endDate {
		startDateTime, err := ConvertDateOnlyStringToDate(startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		return []time.Time{startDateTime}, nil
	}
	layout := "2006-01-02" // The format of the input dates
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %v", err)
	}

	if start.After(end) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	var dates []time.Time
	current := start
	for current.Before(end) || current.Equal(end) {
		dates = append(dates, current)
		current = current.AddDate(0, 0, 1)
	}

	return dates, nil
}

func FindDateDifference(startDate, endDate string) (int, error) {
	if startDate == endDate {
		return 1, nil
	}
	layout := "2006-01-02" // The format of the input dates
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid start date format: %v", err)
	}

	end, err := time.Parse(layout, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid end date format: %v", err)
	}

	if start.After(end) {
		return 0, fmt.Errorf("start date must be before end date")
	}

	difference := end.Sub(start)
	return int(difference.Hours()) / 24, nil
}

func IsWeekend(dateString string) (bool, error) {
	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return false, err
	}

	// Check if the day of the week is Saturday (6) or Sunday (0)
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return true, nil
	}

	return false, nil
}

func AddHoursToTimeString(timeString string, hours int) (string, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return "", err
	}

	newTime := t.Add(time.Duration(hours) * time.Hour)
	newTimeString := newTime.Format("2006-01-02 15:04:05")

	return newTimeString, nil
}

func AddHoursToTime(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

func ConvertDateFormat(inputDateStr string, format string) (string, error) {
	// Parse the input date string in the format "2006-01-02"
	inputDateFormat := "2006-01-02"
	inputDate, err := time.Parse(inputDateFormat, inputDateStr)
	if err != nil {
		return "", err
	}

	// Format the date to the desired output format "02 Jan 2006"
	outputDateFormat := format
	outputDateStr := inputDate.Format(outputDateFormat)

	return outputDateStr, nil
}

func HandleReadableDates(startTime time.Time, endTime time.Time, format string) string {
	startDate := ConvertDateOnlyToString(startTime)
	endDate := ConvertDateOnlyToString(endTime)
	if startDate == endDate {
		return ConvertTimeFormat(startTime, format)
	}
	return fmt.Sprintf("%v to %v", ConvertTimeFormat(startTime, format), ConvertTimeFormat(endTime, format))
}

func ConvertTimeFormat(t time.Time, format string) string {
	formattedTime := t.Format(format)
	return formattedTime
}

const DateDMMYyyy = "2 Jan 2006"

const DateMMDTime = "Jan 2, 3:04 PM"
const DateDMM = "2 Jan"

// ConvertToTimeZone This converts the time to that time zone, but return an error is any found
func ConvertToTimeZone(t time.Time, timeZoneIdentifier string) (time.Time, error) {
	location, err := time.LoadLocation(timeZoneIdentifier)
	if err != nil {
		return t, err
	}

	convertedTime := t.In(location)
	return convertedTime, nil
}

// ConvertToTimeZoneTwo This converts the time to that time zone, but if there is an error we send the original time
func ConvertToTimeZoneTwo(t time.Time, timeZoneIdentifier string) time.Time {
	location, err := time.LoadLocation(timeZoneIdentifier)
	if err != nil {
		return t
	}

	convertedTime := t.In(location)
	return convertedTime
}

func DateByAddOrSubtractDays(dateString string, daysToAddOrSubtract int) (string, error) {
	// Parse the input date string
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateString)
	if err != nil {
		return "", err // Invalid input date string
	}

	// Calculate the new date by adding or subtracting days
	modifiedDate := date.AddDate(0, 0, daysToAddOrSubtract)

	// Format the resulting date as a string
	modifiedDateString := modifiedDate.Format(layout)

	return modifiedDateString, nil
}

func AreDatesInSameMonthAndYear(dateStr1, dateStr2 string) (bool, error) {
	layout := "2006-01-02"

	parsedDate1, err := time.Parse(layout, dateStr1)
	if err != nil {
		return false, err
	}

	parsedDate2, err := time.Parse(layout, dateStr2)
	if err != nil {
		return false, err
	}

	return parsedDate1.Year() == parsedDate2.Year() && parsedDate1.Month() == parsedDate2.Month(), nil
}

