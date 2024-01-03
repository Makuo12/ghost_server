package val


// Validate house rules tag based on the house rules tag
func ValidateHouseRuleTag(tag string) bool {
	for i := 0; i < len(house_rules); i++ {
		if tag == house_rules[i] {
			return true
		}
	}
	return false
}