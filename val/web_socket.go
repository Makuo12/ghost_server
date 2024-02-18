package val

import (
	"regexp"
	"strings"
)

// ROOM TYPE IS FOR WEBSOCKET
var roomTypes = []string{"search_option", "reserve", "search_cal_option", "ex_search_event", "message_listen", "message_unread", "notification_listen", "get_message"}

// Check if it contains room type for websocket
func ContainRoomType(s string) bool {
	for _, room := range roomTypes {
		if room == s {
			return true
		}
	}
	return ValidateMsgRoom(s)
}

const MsgPattern = `^room&[^/]+`
const RoomPattern = `^room&[^/]+`

func CheckPattern(pattern string, input string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(input)
}


func ValidateMsgRoom(room string) bool {
	if CheckPattern(RoomPattern, room) {
		data := strings.Split(room, "&")
		return len(data) == 2
	}
	return false
}