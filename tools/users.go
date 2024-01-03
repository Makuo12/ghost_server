package tools

import "strings"

func HandlePhoneNumber(phone string) string {
	data := ""
	phoneList := strings.Split(phone, "")
	if len(phoneList) < 5 {
		return phone
	} else {
		count := len(phoneList) - 4
		for i := 0; i < count; i++ {
			data += "*"
		}
		for i := count; i < len(phoneList); i++ {
			data += phoneList[i]
		}
	}
	return data
}

func HandleEmail(email string) string {
	data := ""
	emailList := strings.Split(email, "")
	// AfterAt tells us if we have passed @
	afterAt := false
	for i, c := range emailList {
		if c == "@" {
			data += c
			afterAt = true
			continue
		}
		if i == 0 {
			data += c
			continue
		}
		if afterAt {
			data += c
		} else {
			data += "*"
		}
	}
	return data
}
