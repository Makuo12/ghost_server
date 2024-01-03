package val

import "strconv"





func ValidateMoney(balance string) bool {
	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return false
	}
	if bal < 2 {
		return false
	}
	return true
}