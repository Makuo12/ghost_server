package sender

import (
	"context"
	"fmt"
	"strconv"
	"time"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendAdminEmailBrevo(ctx context.Context, cfg *brevo.Configuration, toName string, toEmail string, code string, templateID string, funcName string, appKey string, expire string, userEmail string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"code":      code,
		"year":      fmt.Sprint(time.Now().Year()),
		"useremail": userEmail,
		"expire":    expire,
	}
	dataString := fmt.Sprintf("<html><body><h1>Verify your email address</h1><div>Please verify your email address by entering the six-digit code to continue on Flizzup</div><div>%v</div><footer>Team Flizzup</footer></body></html", code)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendVerifyEmailBrevo(ctx context.Context, cfg *brevo.Configuration, toName string, toEmail string, code string, templateID string, funcName string, appKey string, expire string, adminPhoneNumber string, adminContactEmail string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"code":         code,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
		"expire":       expire,
	}
	dataString := fmt.Sprintf("<html><body><h1>Verify your email address</h1><div>Please verify your email address by entering the six-digit code to continue on Flizzup</div><div>%v</div><footer>Team Flizzup</footer></body></html", code)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendInviteEmailBrevo(ctx context.Context, cfg *brevo.Configuration, toName string, toEmail string, code string, templateID string, funcName string, appKey string, hostName string, options string, expire string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"host":    hostName,
		"code":    code,
		"year":    fmt.Sprint(time.Now().Year()),
		"options": options,
		"expire":  expire,
	}
	dataString := fmt.Sprintf("<html><body><h1>Hosting invitation</h1><h3>%v feels you would make a great host</h3><div>With an extra pair of hands, hosting becomes easier. Help %v host %v</div><div>Invitation code: %v</div><footer>Team Flizzup</footer></body></html", hostName, hostName, options, code)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendCoHostDeactivateBrevo(ctx context.Context, cfg *brevo.Configuration, coHostName string, mainHostEmail string, mainHostName string, optionName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"name":       coHostName,
		"optionname": optionName,
		"year":       fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>Your co-hosting session with %v for %v has ended</h1><div>{{params.name}} just ended co-hosting session with you. If this was a mistake you can just send another invite code to %v's email</div><footer>Team Flizzup</footer></body></html", coHostName, optionName, coHostName)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: mainHostEmail, Name: mainHostName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendReservationRequestBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"header":  header,
		"message": message,
		"expire":  expire,
		"year":    fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, expire)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendAdminReservationRequestBrevo(ctx context.Context, cfg *brevo.Configuration, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, expire string, templateID string, funcName string, appKey string) {
	toEmail := "sylviauwa@gmail.com"
	toName := "Sylvia"
	header := "Reservation request"
	msg := fmt.Sprintf("host_email: %v.\n host_first_name: %v.\n host_last_name: %v.\n charge_id: %v\n. host_public_id: %s\n. guest_email: %v.\n guest_first_name: %v.\n guest_last_name: %v.\n guest_public_id: %v.\n", hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID)
	params := map[string]any{
		"header":  header,
		"message": msg,
		"expire":  expire,
		"year":    fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, msg, expire)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}

	var paramsData map[string]any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "info@flizzup.com",
		},
		Params: paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return

}

//func SendAdminReservationRequestBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
//	//cfg := brevo.NewConfiguration()
//	////Configure API key authorization: api-key
//	//cfg.AddDefaultHeader("api-key", appKey)

//	var paramsData map[string]any = params
//	br := brevo.NewAPIClient(cfg)
//	body := brevo.SendSmtpEmail{
//		HtmlContent: dataString,
//		TemplateId:  int64(tID),
//		To: []brevo.SendSmtpEmailTo{
//			{Email: toEmail, Name: toName},
//		},
//		ReplyTo: &brevo.SendSmtpEmailReplyTo{
//			Name:  "team",
//			Email: "info@flizzup.com",
//		},
//		Params: paramsData,
//	}
//	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
//	if err != nil {
//		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
//		return
//	}
//	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
//	return
//}
