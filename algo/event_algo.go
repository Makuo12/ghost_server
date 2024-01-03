package algo

import (
	"regexp"
)


func HandleEventAlgoDes(des string) map[string]int {
	ratio := make(map[string]int)
	re := regexp.MustCompile(`[?.,\s\d{Zs}":_!;/']+`)
	// lets start
	words := re.Split(des, -1)

	for key, event := range eventDes {
		num := 0
		var rate float64
		for _, w := range words {
			if ContainsString(event, w) {
				num += 1
			}
		}
		// We use 20 because if the length of the words is greater than 20 we just it to have at least 12 of the words to have 100
		if len(event) >= 20 {
			if num >= 12 {
				rate = 100
			} else {
				rate = float64((num * 100) / 12)
			}

		} else {
			// We just need it to have half of the words to be on 100
			value := 0.5 * float64(len(event))
			if float64(num) >= value {
				rate = 100
			} else {
				rate = float64(num*100) / value
			}
		}
		ratio[key] = int((rate * 22) / 100)
	}
	return ratio
}

func HandleEventAlgoName(name string) map[string]int {
	ratio := make(map[string]int)
	re := regexp.MustCompile(`[?.,\s\d{Zs}":_!;/']+`)
	// lets start
	words := re.Split(name, -1)

	for key, event := range eventDes {
		num := 0
		var rate float64
		for _, w := range words {
			if ContainsString(event, w) {
				num += 1
			}
		}
		// We use 20 because if the length of the words is greater than 20 we just it to have at least 12 of the words to have 100
		if len(event) >= 20 {
			if num >= 2 {
				rate = 100
			} else {
				rate = float64((num * 100) / 2)
			}

		} else {
			// We just need it to have half of the words to be on 100
			value := 0.2 * float64(len(event))
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

func HandleEventAlgoHigh(highlights []string) map[string]int {
	ratio := make(map[string]int)

	for key, event := range eventHighlight {
		num := 0
		var rate float64
		for _, a := range highlights {
			if ContainsString(event, a) {
				num += 1
			}
		}
		value := 0.5 * float64(len(event))
		if float64(num) >= value {
			rate = 100
		} else {
			rate = float64(num*100) / value
		}
		ratio[key] = int((rate * 12) / 100)
	}
	return ratio
}

func HandleEventAlgoType(eventTypeData string) map[string]int {
	ratio := make(map[string]int)
	for key, event := range eventType {
		value := 0
		if ContainsString(event, eventTypeData) {

			value = 40
		}
		ratio[key] = value
	}
	return ratio
}

func HandleEventAlgoSubType(eventSubTypeData string) map[string]int {
	ratio := make(map[string]int)
	for key, event := range eventSubType {
		value := 0
		if ContainsString(event, eventSubTypeData) {

			value = 30
		}
		ratio[key] = value
	}
	return ratio
}

