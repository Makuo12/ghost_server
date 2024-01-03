package utils

// This is where we would handle views when going backwards



func LocationViewState(t string) (currentState string, previousState string) {
	switch t {
	case "shortlets":
		currentState, previousState = Amenities, LocationView
		return
	default:
		currentState, previousState = Description, LocationView
		return
	}
}
func DescriptionViewState(t string) (currentState string, previousState string) {
	if t == "shortlets" || t == "lodge" || t == "yatch" {
		currentState, previousState = Name, Description
		return
	}
	return
}
func DescriptionReserveViewState(t string) (currentState string, previousState string) {
	switch t {
	case "shortlets":
		currentState, previousState = Description, Amenities
		return
	case "events":
		currentState, previousState = Description, EventSubType
	default:
		return
	}
	return
}

// It takes in the mainOption parameters as t
//func HighlightViewState(t string) (currentState string, previousState string) {
//	switch t {
//	case "events":
//		currentState, previousState = Publish, Highlight
//		return
//	default:
//		currentState, previousState = Photo, Highlight
//		return
//	}
//}

// It takes in the mainOption parameters as t
// It takes in the OptionType parameters as s
func HighlightReserveViewState(t string, s string) (currentState string, previousState string) {
	//mainOption
	switch t {
	case "events":
		currentState, previousState = Highlight, Name
		return
	case "options":
		currentState, previousState = Highlight, Price
		return
	default:
		return
	}
}

func NameViewState(t string) (currentState string, previousState string) {
	switch t {
	case "shortlets":
		currentState, previousState = Price, Name
		return
	case "events":
		currentState, previousState = Highlight, Name
	default:
		return
	}
	return 
}


func PhotoViewState(t string) (currentState string, previousState string) {
	switch t {
	case "shortlets":
		currentState, previousState = HostQuestion, Photo
		return
	case "events":
		currentState, previousState = HostQuestion, Photo
	default:
		return
	}
	return
}

// It takes in the mainOption parameters as t
// It takes in the OptionType parameters as s
func PublishReverseViewState(t string, s string) (currentState string, previousState string) {
	//mainOption
	switch t {
	case "events":
		currentState, previousState = Publish, Photo
		return
	case "options":
		currentState, previousState = Publish, HostQuestion
	default:
		return
	}
	return
}

//func PhotoReverseViewState(t string) (currentState string, previousState string) {
//	switch t {
//		case "shortlets", "yatch":
//		currentState, previousState = Photo, Photo
//		return
//	default:
//		currentState, previousState = Photo, Photo
//		return
//	}
//}
func LocationReverseViewState(t string) (currentState string, previousState string) {
	switch t {
	case "shortlets":
		currentState, previousState = LocationView, ShortletInfo
		return
	default:
		return
	}
}
