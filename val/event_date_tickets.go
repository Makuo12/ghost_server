package val

var eventTicketTypes = []string{"paid"}

var grades = []string{"general_admission", "early_bird", "vip", "all_inclusive", "premium", "vip_backstage_pass", "gold"}

var eventTicketMainTypes = []string{"ticket", "table"}

func ValidateEventTicketType(s string) (found bool) {
	found = false
	for i := 0; i < len(eventTicketTypes); i++ {
		if s == eventTicketTypes[i] {
			found = true
			return
		}
	}
	return
}

func ValidateEventTicketLevel(s string) (found bool) {
	found = false
	for i := 0; i < len(grades); i++ {
		if s == grades[i] {
			found = true
			return
		}
	}
	return
}

func ValidateEventTicketMainType(s string) (found bool) {
	found = false
	for i := 0; i < len(eventTicketMainTypes); i++ {
		if s == eventTicketMainTypes[i] {
			found = true
			return
		}
	}
	return
}
