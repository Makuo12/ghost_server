package val


var longStayCancel = []string{"long_solid", "long_strict"}

var regStayCancel = []string{"flexible", "moderate", "solid", "strict", "extreme"}

func ValidateLongStayCancel(s string) (found bool) {
	found = false
	for i := 0; i < len(longStayCancel); i++ {
		if s == longStayCancel[i] {
			found = true
			return
		}
	}
	return
}

func ValidateRegStayCancel(s string) (found bool) {
	found = false
	for i := 0; i < len(regStayCancel); i++ {
		if s == regStayCancel[i] {
			found = true
			return
		}
	}
	return
}

type CancelPolicyItem struct {
	Hours    int    `json:"hours"`
	HoursTwo int    `json:"hours_two"`
	Percent  int    `json:"percent"`
	Type     string `json:"type"`
}

type CancelPolicy struct {
	CancelType string             `json:"cancel_type"`
	Items      []CancelPolicyItem `json:"items"`
}

var veryFlexible = []CancelPolicyItem{{24, 0, 100, "standard"}, {12, 0, 50, "standard"}}

var flexible = []CancelPolicyItem{{24, 0, 100, "standard"}}

var veryModerate = []CancelPolicyItem{{120, 0, 100, "standard"}, {48, 0, 50, "standard"}}

var moderate = []CancelPolicyItem{{120, 0, 100, "standard"}}

var solid = []CancelPolicyItem{{720, 0, 100, "standard"}, {48, 240, 100, "hard"}, {120, 0, 50, "standard"}}

var strict = []CancelPolicyItem{{48, 336, 100, "hard"}, {168, 0, 50, "standard"}}

var Policies = []CancelPolicy{
	{
		CancelType: "very_flexible",
		Items:      veryFlexible,
	},
	{
		CancelType: "flexible",
		Items:      flexible,
	},
	{
		CancelType: "very_moderate",
		Items:      veryModerate,
	},
	{
		CancelType: "moderate",
		Items:      moderate,
	},
	{
		CancelType: "solid",
		Items:      solid,
	},
	{
		CancelType: "strict",
		Items:      strict,
	},
}

func GetCancelPolicy(tag string) (CancelPolicy, bool) {
	for _, p := range Policies {
		if p.CancelType == tag {
			return p, true
		}
	}
	return CancelPolicy{}, false
}
