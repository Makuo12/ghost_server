package algo

import (
	"fmt"
	"strings"
)

func HandleAlgoData(data map[string]int) []string {
	result := []string{}
	for key, value := range data {

		result = append(result, fmt.Sprintf("%v&%v", key, value))
	}
	return result
}

func ContainsString(a []string, s string) bool {
	for i := 0; i < len(a); i++ {
		if strings.EqualFold(a[i], s) {
			return true
		}
	}
	return false
}


