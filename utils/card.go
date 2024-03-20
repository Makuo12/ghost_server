package utils

import "regexp"

func MatchCardType(input string) string {
	// Define regular expression pattern
	pattern := `mastercard|visa|verve`

	// Compile regular expression pattern
	regExp := regexp.MustCompile(pattern)

	// Find the first match in the input string
	match := regExp.FindString(input)

	// Return the matched string if found
	if match != "" {
		return match
	}

	// Return empty string if no match found
	return input
}
