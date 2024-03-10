package val

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	isValidUsername    = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidName        = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
	isValidPhoneNumber = regexp.MustCompile(`\d+\s\d+`).MatchString
)

var userProfileTypes = []string{"work", "language", "bio"}

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidateName(value string) bool {
	//if err := ValidateString(value, 3, 100); err != nil {
	//	return false
	//}
	//if !isValidName(value) {
	//	return false
	//}
	return true
}

func hasSQLInjection(input string) bool {
	// Regular expression to match common SQL injection patterns
	sqlRegex := regexp.MustCompile(`(?i)\b(union|select|from|where|and|or|insert|update|delete|truncate|drop|alter|create|database)\b`)

	// Convert input to lowercase for case-insensitive matching
	lowercaseInput := strings.ToLower(input)

	// Check if the input matches the SQL injection pattern
	return sqlRegex.MatchString(lowercaseInput)
}

func VerifyPassword(s string) bool {
	//var hasNumber, hasUpperCase, hasLowercase, hasSpecial bool
	//if len(s) < 7 || hasSQLInjection(s) {
	//	return false
	//}

	//for _, c := range s {
	//	switch {
	//	case unicode.IsNumber(c):
	//		hasNumber = true
	//	case unicode.IsUpper(c):
	//		hasUpperCase = true
	//	case unicode.IsLower(c):
	//		hasLowercase = true
	//	case c == '#' || c == '|':
	//		return false
	//	case unicode.IsPunct(c) || unicode.IsSymbol(c):
	//		hasSpecial = true

	//	}
	//}
	//return hasNumber && hasUpperCase && hasLowercase && hasSpecial
	return true
}

func ValidateEmail(value string) bool {
	if err := ValidateString(value, 3, 200); err != nil {
		return false
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return false
	}
	return true
}

func ValidatePhoneNumber(n string) bool {
	return isValidPhoneNumber(n)
}

func ValidateUserProfileType(s string) bool {
	for i := 0; i < len(userProfileTypes); i++ {
		if s == userProfileTypes[i] {
			return true
		}
	}
	return false
}
