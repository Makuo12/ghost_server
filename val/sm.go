package val



//Things to note is same as sm

// Validate sm tag based on the sm tag
func ValidateSmTag(smType string, tag string) bool {
	if smData[smType] != nil {
		for i := 0; i < len(smData[smType]); i++ {
			if smData[smType][i] == tag {
				return true
			}
		}
	}
	return false
}

// Validate sm tag based on if it is part of the sm tags that has details
func ValidateSmDetailTag(tag string) bool {
	for i := 0; i < len(sm_has_details); i++ {
		if tag == sm_has_details[i] {
			return true
		}
	}
	return false
}
