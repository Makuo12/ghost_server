package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/sender"
	"github.com/makuo12/ghost_server/utils"
)

//func SendEmailVerifyCode(server *Server, toEmail string, toName string, usernameString string, funcName string) (err error) {
//	code := utils.RandomNumber(5)
//	// We want to send the email to msg91
//	// We wouldn't set it on right now
//	err = sender.SendEmailVerifyCode(toName, toEmail, code, server.config.AppName, server.config.AppEmail, server.config.AppEmailDomain, server.config.EmailTemplate, funcName, server.config.Msg91Key)
//	if err != nil {
//		return
//	}
//	log.Println("Flizzup code is: ", code)
//	userDetails := fmt.Sprintf("%v&%v&%v&%v&%v", "email", code, usernameString, toEmail, "false")
//	err = RedisClient.Set(RedisContext, usernameString, userDetails, time.Hour*4).Err()
//	if err != nil {
//		log.Printf("FuncName: %v, SendEmailVerifyCode RedisClient.HSe %v err:%v\n", funcName, usernameString, err.Error())
//		err = fmt.Errorf("code not store code, try again")
//		return
//	}
//	return
//}

//func SendEmailInvitationCode(server *Server, toEmail string, toName string, hostNameOption string, mainHostName string, mainOption string, funcName string, coID uuid.UUID) (err error) {
//	code := utils.RandomNumber(5)
//	// We want to send the email to msg91
//	td := time.Hour * 28
//	t := time.Now().Add(td)
//	timeString := tools.ConvertTimeFormat(t, tools.DateMMDTime)
//	log.Println("Flizzup cohost code is: ", code, "\n time is ", timeString)
//	err = sender.SendEmailInvitationCode(toName, toEmail, hostNameOption, mainHostName, mainOption, timeString, code, server.config.AppName, server.config.AppEmail, server.config.AppEmailDomain, server.config.InviteTemplate, funcName, server.config.Msg91Key)
//	if err != nil {
//		return
//	}
//	err = RedisClient.Set(RedisContext, string(code), tools.UuidToString(coID), td).Err()
//	if err != nil {
//		log.Printf("FuncName: %v, SendEmailInvitationCode RedisClient.Set err:%v\n", funcName, err.Error())
//		err = fmt.Errorf("code not store code, try sending invitation again")
//		return
//	}
//	return
//}

//func SendCustomEmail(server *Server, toEmail string, toName string, header string, topHeader string, body string, funcName string) (err error) {
//	err = sender.SendCustomEmail(toName, toEmail, header, topHeader, body, server.config.AppName, server.config.AppEmail, server.config.AppEmailDomain, server.config.EmailTemplate, funcName, server.config.Msg91Key)
//	if err != nil {
//		return
//	}
//	return
//}

func SendSmsOtp(server *Server, toPhone string, usernameString string, dial_number_code int32, dial_country string, funcName string) (err error) {
	code := utils.RandomNumber(5)
	dial_code_number := strconv.FormatInt(int64(dial_number_code), 10)
	data := fmt.Sprintf("%v&%v&%v&%v&%v&%v&%v", "phone", code, toPhone, usernameString, dial_code_number, dial_country, "false")
	split := strings.Split(toPhone, "_")
	if len(split) != 2 {
		err = fmt.Errorf("phone in wrong format")
	}
	sendNumber := fmt.Sprintf("%v%v", split[0], split[1])
	err = sender.SendSmsOtp(funcName, server.config.Msg91Key, server.config.SmsOtpTemplate, sendNumber, code)
	if err != nil {
		return
	}
	log.Println("Flizzup phone code is: ", code)
	err = RedisClient.Set(RedisContext, usernameString, data, time.Hour*4).Err()
	if err != nil {
		log.Printf("SendSmsOtp SAdd(RedisContext, RedisClient.Set %v err:%v\n", usernameString, err.Error())
	}
	log.Printf("Your Flizzup verification code is: %s\n", code)
	return
}

func SendSmsOtpWithName(server *Server, toPhone string, usernameString string, dial_number_number string, dial_country string, firstName string, lastName string, funcName string) (err error) {
	code := utils.RandomNumber(5)
	data := fmt.Sprintf("%v&%v&%v&%v&%v&%v&%v&%v&%v", "phone", code, toPhone, usernameString, dial_number_number, dial_country, "false", firstName, lastName)
	split := strings.Split(toPhone, "_")
	if len(split) != 2 {
		err = fmt.Errorf("phone in wrong format")
	}
	sendNumber := fmt.Sprintf("%v%v", split[0], split[1])
	err = sender.SendSmsOtp(funcName, server.config.Msg91Key, server.config.SmsOtpTemplate, sendNumber, code)
	if err != nil {
		return
	}
	log.Println("Flizzup phone code is: ", code)
	err = RedisClient.Set(RedisContext, usernameString, data, time.Hour*4).Err()
	if err != nil {
		log.Printf("SendSmsOtp SAdd(RedisContext, RedisClient.Set %v err:%v\n", usernameString, err.Error())
	}
	log.Printf("Your Flizzup verification code is: %s\n", code)
	return
}

func StoreNumber(server *Server, toPhone string, usernameString string, dial_number_code int32, dial_country string, funcName string) (err error) {
	dial_code_number := strconv.FormatInt(int64(dial_number_code), 10)
	data := fmt.Sprintf("%v&%v&%v&%v&%v&%v&%v", "phone", "no_code", toPhone, usernameString, dial_code_number, dial_country, "false")
	err = RedisClient.Set(RedisContext, usernameString, data, time.Hour*4).Err()
	if err != nil {
		log.Printf("SendSmsOtp SAdd(RedisContext, RedisClient.Set %v err:%v\n", usernameString, err.Error())
	}
	return
}
