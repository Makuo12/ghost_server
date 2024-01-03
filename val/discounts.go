package val

var optionDiscounts = []string{"weekly_discount", "monthly_discount", "week_8_discount", "week_12_discount", "week_16_discount"}

type DiscountItem struct {
	Type   string `json:"type"`
	Number int    `json:"number"`
}

var optionDiscountData = []DiscountItem{
	{
		Type:   "weekly_discount",
		Number: 7,
	},
	{
		Type:   "monthly_discount",
		Number: 28,
	},
	{
		Type:   "week_8_discount",
		Number: 56,
	},
	{
		Type:   "week_12_discount",
		Number: 84,
	},
	{
		Type:   "week_16_discount",
		Number: 112,
	},
}

// This takes in a discount type and returns the discount number
func GetOptionDiscountNumber(t string) int {
	for _, v := range optionDiscountData {
		if v.Type == t {
			return v.Number
		}
	}
	return 0
}


func ValidateOptionDiscount(s string) bool {
	for i:=0; i<len(optionDiscounts); i++ {
		if optionDiscounts[i] == s {
			return true
		}
	}
	return false
}

