package sender

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func SendSmsOtp(funcName, key, templateID, phone, otp string) (err error) {
	url := fmt.Sprintf("https://control.msg91.com/api/v5/otp?template_id=%v&mobile=%v&otp=%v", templateID, phone, otp)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode SendSmsOtp error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send verification code")
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode http.DefaultClient.Do( error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send verification code")
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode http.DefaultClient.Do( error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send verification code")
		return
	}
	if res.StatusCode >= 400 {
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	log.Println(string(body))
	return

}
