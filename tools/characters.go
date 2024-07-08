package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ConvertStringToBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}
}

func ConvertInt64ToString(num int64) string {
	return strconv.Itoa(int(num))
}

func ConvertInt32ToString(num int32) string {
	return strconv.Itoa(int(num))
}

func ConvertStringToInt32(num string) (int32, error) {
	index, err := strconv.Atoi(strings.TrimSpace(num))
	return int32(index), err
}

func ConvertStringToInt64(num string) (int64, error) {
	index, err := strconv.Atoi(strings.TrimSpace(num))
	return int64(index), err
}
func ConvertBoolToString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

// Converts float to string with two decimal places
func ConvertFloatToString(num float64) string {
	result := strconv.FormatFloat(num, 'f', 2, 64)
	return result
}

func ConvertLocationStringToFloat(latitudeStr string, precision int) float64 {
	// Parse the latitude string to a float64 with the specified precision
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		log.Printf("Error at ConvertStringToFloat %v\n", err)
		return 0.00
	}
	return latitude
}

func ConvertFloatToLocationString(latitude float64, precision int) string {
	// Use fmt.Sprintf to format the float with the specified precision
	formatString := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(formatString, latitude)
}

func ConvertStringToFloat(num string) float64 {
	result, err := strconv.ParseFloat(num, 64)
	if err != nil {
		log.Printf("Error at ConvertStringToFloat %v\n", err)
		return 0.00
	}
	return result
}

func ContainsString(a []string, s string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == s {
			return true
		}
	}
	return false
}

// Remove the string and returns a slice without the string
func RemoveString(a []string, s string) []string {
	var data []string
	for i := 0; i < len(a); i++ {
		if a[i] != s {
			data = append(data, a[i])
		}
	}
	return data
}

func ExtractNumberFromString(input string) (int, error) {
	// Define a regular expression to match numbers in the input string
	re := regexp.MustCompile(`\d+`)

	// Find all occurrences of numbers in the string
	numbers := re.FindAllString(input, -1)

	// If no number is found, return an error
	if len(numbers) == 0 {
		return 0, fmt.Errorf("no number found in the string")
	}

	// Convert the first number to an integer and return it
	number, err := strconv.Atoi(numbers[0])
	if err != nil {
		return 0, fmt.Errorf("failed to convert number to integer: %v", err)
	}

	return number, nil
}

func CapitalizeFirstCharacter(input string) string {
	if len(input) == 0 {
		return input
	}

	// Convert string to rune slice
	runes := []rune(input)

	// Capitalize the first rune if it's a letter
	if unicode.IsLetter(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}

	// Convert rune slice back to string
	return string(runes)
}

// Convert metadata to string
func ConvertAnyToString(metadata any) (string, error) {
	if metadata == nil {
		return "", nil
	}

	// Marshal the metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}

	return string(metadataJSON), nil
}

func RemoveDuplicates(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}

	for _, item := range slice {
		if _, found := seen[item]; !found {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func GetImageItem(image string) (path string, url string) {
	log.Println("images: ", image)
	split := strings.Split(image, "*")
	log.Println("images split: ", split)
	if len(split) == 2 {
		path = split[0]
		url = split[1]
		return
	}
	return
}

func GetImageListItem(images []string) ([]string, []string) {
	myPaths := []string{}
	myUrls := []string{}
	for _, image := range images {
		split := strings.Split(image, "*")
		if len(split) == 2 {
			myPaths = append(myPaths, split[0])
			myUrls = append(myUrls, split[0])
		}
	}
	return myPaths, myUrls
}
