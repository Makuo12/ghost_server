package val

import (
	"regexp"
)

var (
	isValidDateOnly = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`).MatchString
	isValidTimeOnly = regexp.MustCompile(`\d\d:\d\d`).MatchString
)

var months = []string{"January","February","March","April","May","June","July","August","September","October","November","December"}


func ValidateTimeOnly(s string) bool {
	return isValidTimeOnly(s)
}


func ValidateDateOnly(s string) bool {
	return isValidDateOnly(s)
}


func ValidateMonth(s string) bool {
	for i:=0; i<len(months); i++ {
		if months[i] == s {
			return true
		}
	}
	return false
}