package val

var ShortletSpace = []string{
	"full_place",
	"private_room",
	"shared_room",
}

var GuestTypes = []string{"adult", "children", "infant", "pet"}

var Shortlets = []string{"flat_apartment", "home", "bed_and_breakfast", "guest_house", "barn", "yacht", "tiny_home", "nature_lodge", "tower", "equipped_apartment", "bungalow", "castle", "minsu", "container", "cave", "vacation_home", "cottage", "casa", "tree_house", "windmill", "house_boat", "motorhome", "light_house", "lighthouse", "dome", "plane", "bus", "loft", "boy_quarter", "cabin"}

var OptionExtraInfo = []string{"help_manual", "direction"}

var GuestAreas = []string{"bedroom", "full_bathroom", "half_bathroom", "kitchen", "half_kitchen", "living_room", "dining_area", "office", "back_garden", "patio", "pool", "gym", "hot_tub"}

var PropertySizeUnits = []string{"Sq ft", "Sq m"}

var DesTypes = []string{"des", "space_des", "guest_access_des", "interact_with_guests_des", "other_des", "neighborhood_des", "get_around_des"}

func ValidateOptionType(s string) bool {
	for i := 0; i < len(Options); i++ {
		if s == Options[i] {
			return true
		}
	}
	return false
}

func ValidateShortletType(s string) bool {
	for i := 0; i < len(Shortlets); i++ {
		if s == Shortlets[i] {
			return true
		}
	}
	return false
}

func ValidateOptionExtraInfo(s string) bool {
	for i := 0; i < len(OptionExtraInfo); i++ {
		if s == OptionExtraInfo[i] {
			return true
		}
	}
	return false
}

func ValidateShortletSpace(s string) bool {
	for i := 0; i < len(ShortletSpace); i++ {
		if s == ShortletSpace[i] {
			return true
		}
	}
	return false
}

func ValidateDesTypes(s string) bool {
	for i := 0; i < len(DesTypes); i++ {
		if s == DesTypes[i] {
			return true
		}
	}
	return false
}

func ValidatePropertyUnit(u string) bool {
	for _, p := range PropertySizeUnits {
		if p == u {
			return true
		}
	}
	return false
}

func ContainsBedroomAndNumGuest(s []string) (bedroomFound bool, numGuest int) {
	bedroomFound = false
	numGuest = 0
	for _, a := range s {
		if a == "bedroom" {
			bedroomFound = true
		}
		if a == "guest" {
			numGuest += 1
		}
	}
	return bedroomFound, numGuest
}


func ValidateGuestOptionTypes(guests []string) bool {
	for _, g := range guests {
		exist := false
		for _, s := range GuestTypes {
			if s == g {
				exist = true
			}
		}
		if !exist {
			return false
		}
	}
	return true
}