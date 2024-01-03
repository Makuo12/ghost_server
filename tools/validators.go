package tools

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"strconv"
)

var am_has_details = []string{"dedicated_workspace", "air_condition", "dryer", "heating", "hot_tub", "tv", "washing_machine", "body_soap", "conditioner", "shampoo", "shower_gel", "clothes_storage", "exercise_equipment", "game_console", "piano", "sound_system", "basketball_court", "tennis_court", "football_court", "children_books_and_toy", "high_chair", "generator", "electricity", "dedicated_workspace", "indoor_fireplace", "coffee_maker", "oven", "refrigerator", "stove", "beach_access", "resort_access", "ski_in_out", "garden", "bbq_grill", "outdoor_kitchen", "free_parking_on_premises", "paid_parking_on_premises", "paid_parking_off_premises", "pool", "sauna"}

func CheckStringIsFloat(num string) bool {
	_, err := strconv.ParseFloat(num, 64)
	return err == nil
}

func CheckAmHasDetail(am string) bool {
	var hasDetails = false
	for _, v := range am_has_details {
		if v == am {
			hasDetails = true
		}
	}
	return hasDetails
}

func ValidateIntLessThanZero(val int) int {
	if val < 0 {
		return 0
	} else {
		return val
	}
}

func IsValidSignature(data []byte, signature, secretKey string, hashFor string) bool {
	mac := hmac.New(sha512.New, []byte(secretKey))
	_, err := mac.Write(data)
	if err != nil {
		log.Printf("error at IsValidSignature for %v err: %v\n", hashFor, err)
		return false
	}
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedMAC), []byte(signature))
}

func HandleHMAC(reqBody []byte, secret string) (hash string, err error) {
	// Generate HMAC
	hash = generateHMAC(secret, string(reqBody))
	return
}

func generateHMAC(secret, data string) string {
	key := []byte(secret)
	message := []byte(data)

	h := hmac.New(sha512.New, key)
	h.Write(message)
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}


