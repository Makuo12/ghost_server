package tools

import (
	"errors"
	"fmt"
	"time"
)

func ServerListIsEmpty(a []string) bool {
	if len(a) == 0 {
		return true
	} else if len(a) == 1 {
		if a[0] == "" || a[0] == "none" {
			return true
		}
	}
	return false
}



// Remove none and empty values from a list
func HandleListReq(a []string) []string {
	var arr []string
	for i := 0; i < len(a); i++ {
		if a[i] != "none" && len(a[i]) > 0 {
			arr = append(arr, a[i])
		}
	}
	return arr
}

func ServerStringEmpty(s string) bool {
	return len(s) == 0 || s == "none"
}

func ServerDoubleEmpty(s string) bool {
	return len(s) == 0 || s == "none" || s == "0.0"
}

func HandleString(s string) string {
	if ServerStringEmpty(s) {
		return ""
	}
	return s
}

// If string is empty it returns null
func HandleStringTwo(s string) string {
	if ServerStringEmpty(s) {
		return "none"
	}
	return s
}

// Returns a new list that would have none if it was empty
func ServerListToDB(a []string) []string {
	if ServerListIsEmpty(a) {
		return []string{"none"}
	}
	return HandleListReq(a)
}

// Handles confirm and err variables
func HandleConfirmError(err error, confirm bool, msg string) error {
	if err != nil {
		return err
	}
	return fmt.Errorf(msg)
}

// This would check if the status of the option is staged if staged then we want to change it to list
func HandleOptionStatus(s string) string {
	if s == "staged" {
		return "list"
	}
	return s
}

func SnoozeDatesGood(startDate time.Time, endDate time.Time) error {
	if startDate.After(endDate) {
		return errors.New("start snooze date must be before your end snooze date")
	} else if startDate.Before(time.Now().Add(time.Hour * 4)) {
		return errors.New("start snooze date must be 4 hours after the current date")
	}
	return nil
}

var FakeDate = "1777-12-07"



func IsInList(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
