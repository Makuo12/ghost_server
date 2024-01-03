package algo

import (
	"regexp"
)



func HandleOptionAlgoDes(des string) map[string]int {
	ratio := make(map[string]int)
	re := regexp.MustCompile(`[?.,\s\d{Zs}":_!;/']+`)
	// lets start
	words := re.Split(des, -1)

	for key, option := range optionDes {
		num := 0
		var rate float64
		for _, w := range words {
			if ContainsString(option, w) {
				num += 1
			}
		}
		// We use 20 because if the length of the words is greater than 20 we just it to have at least 12 of the words to have 100
		if len(option) >= 20 {
			if num >= 12 {
				rate = 100
			} else {
				rate = float64((num * 100) / 12)
			}

		} else {
			// We just need it to have half of the words to be on 100
			value := 0.5 * float64(len(option))
			if float64(num) >= value {
				rate = 100
			} else {
				rate = float64(num*100) / value
			}
		}
		ratio[key] = int((rate * 18) / 100)
	}
	return ratio
}

func HandleOptionAlgoName(name string) map[string]int {
	ratio := make(map[string]int)
	re := regexp.MustCompile(`[?.,\s\d{Zs}":_!;/']+`)
	// lets start
	words := re.Split(name, -1)

	for key, option := range optionDes {
		num := 0
		var rate float64
		for _, w := range words {
			if ContainsString(option, w) {
				num += 1
			}
		}
		// We use 20 because if the length of the words is greater than 20 we just it to have at least 12 of the words to have 100
		if len(option) >= 20 {
			if num >= 2 {
				rate = 100
			} else {
				rate = float64((num * 100) / 2)
			}

		} else {
			// We just need it to have half of the words to be on 100
			value := 0.2 * float64(len(option))
			if float64(num) >= value {
				rate = 100
			} else {
				rate = float64(num*100) / value
			}
		}
		ratio[key] = int((rate * 6) / 100)
	}

	return ratio
}

func HandleOptionAlgoAmenities(amenities []string) map[string]int {
	ratio := make(map[string]int)

	for key, option := range optionAmenity {
		num := 0
		var rate float64
		for _, a := range amenities {
			if ContainsString(option, a) {
				num += 1
			}
		}
		value := 0.5 * float64(len(option))
		if float64(num) >= value {
			rate = 100
		} else {
			rate = float64(num*100) / value
		}
		ratio[key] = int((rate * 12) / 100)
	}
	return ratio
}

func handleSpaceAreaData(spaces []string) []string {
	spacesData := []string{}
	// First we want to convert it to a slice
	// create a map with bool values
	set := make(map[string]bool)
	// loop through the slice and add each element as a key to the map
	for _, v := range spaces {
		set[v] = true
	}

	for space := range set {
		switch space {
		case "full_bathroom", "half_bathroom":
			spacesData = append(spacesData, "bathroom")
		case "full_kitchen", "half_kitchen":
			spacesData = append(spacesData, "kitchen")

		default:
			spacesData = append(spacesData, space)

		}
	}
	return spacesData

}

func HandleOptionAlgoSpaceAreas(spaceAreas []string) map[string]int {
	ratio := make(map[string]int)
	spaceAreas = handleSpaceAreaData(spaceAreas)
	for key, option := range optionSpaceArea {
		num := 0
		var rate float64
		for _, a := range spaceAreas {
			if ContainsString(option, a) {
				num += 1
			}
		}
		value := 0.8 * float64(len(option))
		if float64(num) >= value {
			rate = 100
		} else {
			rate = float64(num*100) / value
		}
		ratio[key] = int((rate * 6) / 100)
	}
	return ratio
}

func HandleOptionAlgoHigh(highlights []string) map[string]int {
	ratio := make(map[string]int)

	for key, option := range optionHighlight {
		num := 0
		var rate float64
		for _, a := range highlights {
			if ContainsString(option, a) {
				num += 1
			}
		}
		value := 0.5 * float64(len(option))
		if float64(num) >= value {
			rate = 100
		} else {
			rate = float64(num*100) / value
		}
		ratio[key] = int((rate * 12) / 100)
	}
	return ratio
}

func HandleOptionAlgoType(shortletType string) map[string]int {
	ratio := make(map[string]int)
	for key, option := range optionType {
		value := 0
		if ContainsString(option, shortletType) {

			value = 40
		}
		ratio[key] = value
	}
	return ratio
}

func HandleOptionAlgoSpaceType(spaceType string) map[string]int {
	ratio := make(map[string]int)
	for key, option := range optionSpaceType {
		value := 0
		if ContainsString(option, spaceType) {

			value = 20
		}
		ratio[key] = value
	}
	return ratio
}


