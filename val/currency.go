package val

// Constants for all supported currencies
const (
	USD = "USD"
	NGN = "NGN"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, NGN, CAD:
		return true
	}
	return false
}