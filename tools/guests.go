package tools

// This function takes in a list of elements then return a map showing the amount of times each element show up
func HandleListCount(s []string)  map[string]int {
	data := make(map[string]int)
	if len(s) != 0 {
		for _, e := range s {
			data[e] += 1
		}
	}
	return data
}


