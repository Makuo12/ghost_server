package val

var EventLocationTypes = []string{"venue", "to_be_announced"}

var EventDateTypes = []string{"recurring_event", "single_event"}

var Events = []string{
	"concert",
	"pray",
	"gather",
	"talk",
	"comedy",
	"pool_party",
	"movie",
	"sport",
}
var Options = []string{
	"shortlets",
	"events",
}


func ValidateLocationType(s string) (found bool){
	found = false
	for i:=0;i<len(EventLocationTypes);i++{
		if s == EventLocationTypes[i] {
			found = true
			return
		}
	}
	return
}

func ValidateEventDateType(s string) (found bool){
	found = false
	for i:=0;i<len(EventDateTypes);i++{
		if s == EventDateTypes[i] {
			found = true
			return
		}
	}
	return
}

func ValidateEventType(s string) bool {

	for i := 0; i < len(Events); i++ {
		if s == Events[i] {
			return true
		}
	}
	return false
}