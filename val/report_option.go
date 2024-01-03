package val

var reportOption = []string{"inaccurate", "not_real", "offensive", "something_else"}

var reportScamOption = []string{"money_transfer", "shared_contact", "advertising", "duplicate", "misleading"}

var reportOffensiveOption = []string{"discriminatory", "inappropriate", "advertising"}

func ValidateReportOption(tag string) bool {
	for i := 0; i < len(reportOption); i++ {
		if tag == reportOption[i] {
			return true
		}
	}
	return false
}

func ValidateReportScamOption(tag string) bool {
	for i := 0; i < len(reportScamOption); i++ {
		if tag == reportScamOption[i] {
			return true
		}
	}
	return false
}

func ValidateReportOffensiveOption(tag string) bool {
	for i := 0; i < len(reportOffensiveOption); i++ {
		if tag == reportOffensiveOption[i] {
			return true
		}
	}
	return false
}
