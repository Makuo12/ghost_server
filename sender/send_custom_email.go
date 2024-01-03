package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendCustomEmail(toName string, toEmail string, header string, topHeader string, body string, appName string, appEmail string, appEmailDomain string, templateID string, funcName string, appKey string) (err error) {

	url := "https://control.msg91.com/api/v5/email/send"

	from := From{
		Name:  appName,
		Email: appEmail,
	}
	variable := VariableCustom{
		CompanyName: appName,
		Header:      header,
		TopHeader:   topHeader,
		Body:        body,
		Year:        string(rune(time.Now().Year())),
	}
	to := []To{{
		Name:  toName,
		Email: toEmail,
	}}

	rec := []RecipientCustom{{
		To:        to,
		Variables: variable,
	}}
	email := EmailCustomMessage{
		Recipients: rec,
		From:       from,
		Domain:     appEmailDomain,
		TemplateID: templateID,
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(email)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode SendCustomEmail error %v\n", funcName, err.Error())
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("Error at funcName: %v, http.NewRequest(POST SendCustomEmail error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", appKey)
	clientSide := &http.Client{}
	res, err := clientSide.Do(req)
	if err != nil {
		log.Printf("Error at funcName: %v,clientSide.Do(req) SendCustomEmail error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	return
}
