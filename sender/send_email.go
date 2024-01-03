package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendEmailVerifyCode(toName string, toEmail string, code string, appName string, appEmail string, appEmailDomain string, templateID string, funcName string, appKey string) (err error) {

	url := "https://control.msg91.com/api/v5/email/send"

	from := From{
		Name:  appName,
		Email: appEmail,
	}
	variable := Variable{
		CompanyName: appName,
		Code:        code,
		Year:        fmt.Sprint(time.Now().Year()),
	}
	to := []To{{
		Name:  toName,
		Email: toEmail,
	}}

	rec := []Recipient{{
		To:        to,
		Variables: variable,
	}}
	email := EmailMessage{
		Recipients: rec,
		From:       from,
		Domain:     appEmailDomain,
		TemplateID: templateID,
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(email)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode SendEmailVerifyCode error %v\n", funcName, err.Error())
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("Error at funcName: %v, http.NewRequest(POST SendEmailVerifyCode error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", appKey)
	clientSide := &http.Client{}
	res, err := clientSide.Do(req)
	if err != nil {
		log.Printf("Error at funcName: %v,clientSide.Do(req) SendEmailVerifyCode error %v\n", funcName, err.Error())
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

func SendEmailInvitationCode(toName string, toEmail string, hostNameOption string, mainHost string, mainOption string, expireTime string, code string, appName string, appEmail string, appEmailDomain string, templateID string, funcName string, appKey string) (err error) {

	url := "https://control.msg91.com/api/v5/email/send"

	from := From{
		Name:  appName,
		Email: appEmail,
	}
	variable := VariableInvitation{
		CompanyName:    appName,
		Code:           code,
		HostNameOption: hostNameOption,
		MainHost:       mainHost,
		MainOption:     mainOption,
		Time:           expireTime,
		Year:           fmt.Sprint(time.Now().Year()),
	}
	to := []To{{
		Name:  toName,
		Email: toEmail,
	}}

	rec := []RecipientInvitation{{
		To:        to,
		Variables: variable,
	}}
	email := EmailInvitationMessage{
		Recipients: rec,
		From:       from,
		Domain:     appEmailDomain,
		TemplateID: templateID,
	}
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(email)
	if err != nil {
		log.Printf("Error at funcName: %v, json.NewEncoder(buf).Encode SendEmailInvitationCode error %v\n", funcName, err.Error())
		return
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		log.Printf("Error at funcName: %v, http.NewRequest(POST SendEmailInvitationCode error %v\n", funcName, err.Error())
		err = fmt.Errorf("could not send the verification email, try again")
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authkey", appKey)
	clientSide := &http.Client{}
	res, err := clientSide.Do(req)
	if err != nil {
		log.Printf("Error at funcName: %v,clientSide.Do(req) SendEmailInvitationCode error %v\n", funcName, err.Error())
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
