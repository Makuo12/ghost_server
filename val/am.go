package val

// Validate am tag based on the am tag
func ValidateAmTag(amType string, tag string) bool {
	if amData[amType] != nil {
		for i := 0; i < len(amData[amType]); i++ {
			if amData[amType][i] == tag {
				return true
			}
		}
	}
	return false
}
// Validate am tag based on if it is part of the am tags that has details
func ValidateAmDetailTag(tag string) bool {
	for i := 0; i < len(am_has_details); i++ {
		if tag == am_has_details[i] {
			return true
		}
	}
	return false
}
