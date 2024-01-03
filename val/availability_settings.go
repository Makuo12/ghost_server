package val


var AdvanceNotice = []string{"same day", "at least 1 day's notice", "at least 2 day's notice", "at least 3 day's notice", "at least 5 day's notice", "at least 7 day's notice"}

var PrepareTime = []string{"none", "for 1 night before & after", "for 2 nights before and after"}

var AvailableWindow = []string{"all future dates", "12 months in advance", "9 months in advance", "6 months in advance", "3 months in advance", "dates unavailable by default"}

// ValidateAvailability s is the input the user entered and t is the type 
func ValidateAvailability(s string, t string) (found bool) {
	found = false
	switch t {
	case "advance_notice":
		for i := 0; i < len(AdvanceNotice); i++ {
			if s == AdvanceNotice[i] {
				found = true
				return
			}
		}
	case "preparation_time":
		for i := 0; i < len(PrepareTime); i++ {
			if s == PrepareTime[i] {
				found = true
				return
			}
		}
	case "availability_window":
		for i := 0; i < len(AvailableWindow); i++ {
			if s == AvailableWindow[i] {
				found = true
				return
			}
		}
	}
	return
}


