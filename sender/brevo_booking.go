package sender

import (
	"context"
	"fmt"
	"strconv"
	"time"

	brevo "github.com/getbrevo/brevo-go/lib"
)

//func SendReservationRequestDisapprovedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
//	//cfg := brevo.NewConfiguration()
//	////Configure API key authorization: api-key
//	//cfg.AddDefaultHeader("api-key", appKey)
//	params := map[string]any{
//		"header":  header,
//		"message": message,
//		"year":    fmt.Sprint(time.Now().Year()),
//	}
//	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, expire)
//	tID, err := strconv.Atoi(templateID)
//	if err != nil {
//		// Handle error if conversion fails
//		fmt.Println("Error:", err)
//		return
//	}
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
//	fmt.Println("SendReservationRequestDisapprovedBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
//	return
//}

func SendAdminReservationRequestDisapprovedBrevo(ctx context.Context, cfg *brevo.Configuration, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, expire string, templateID string, funcName string, appKey string) {
	toEmail := "sylviauwa@gmail.com"
	toName := "Sylvia"
	header := "Reservation request disapproved"
	msg := fmt.Sprintf("host_email: %v.\n host_first_name: %v.\n host_last_name: %v.\n charge_id: %v\n. host_public_id: %s\n. guest_email: %v.\n guest_first_name: %v.\n guest_last_name: %v.\n guest_public_id: %v.\n", hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID)
	params := map[string]any{
		"header":  header,
		"message": msg,
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

func SendPaymentSuccessBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
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
	fmt.Println("SendPaymentSuccessBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendAdminPaymentSuccessBrevo(ctx context.Context, cfg *brevo.Configuration, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, expire string, templateID string, funcName string, appKey string) {
	toEmail := "sylviauwa@gmail.com"
	toName := "Sylvia"
	header := "Payment successful"
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

//func SendPaymentFailedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
//	//cfg := brevo.NewConfiguration()
//	////Configure API key authorization: api-key
//	//cfg.AddDefaultHeader("api-key", appKey)
//	params := map[string]any{
//		"header":  header,
//		"message": message,
//		"expire":  expire,
//		"year":    fmt.Sprint(time.Now().Year()),
//	}
//	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, expire)
//	tID, err := strconv.Atoi(templateID)
//	if err != nil {
//		// Handle error if conversion fails
//		fmt.Println("Error:", err)
//		return
//	}
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
//	fmt.Println("SendPaymentFailedBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
//	return
//}

func SendAdminPaymentFailedBrevo(ctx context.Context, cfg *brevo.Configuration, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, expire string, templateID string, funcName string, appKey string) {
	toEmail := "sylviauwa@gmail.com"
	toName := "Sylvia"
	header := "Payment Failed"
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

//func SendReservationRequestApprovedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, expire string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
//	//cfg := brevo.NewConfiguration()
//	////Configure API key authorization: api-key
//	//cfg.AddDefaultHeader("api-key", appKey)
//	params := map[string]any{
//		"header":  header,
//		"message": message,
//		"expire":  expire,
//		"year":    fmt.Sprint(time.Now().Year()),
//	}
//	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, expire)
//	tID, err := strconv.Atoi(templateID)
//	if err != nil {
//		// Handle error if conversion fails
//		fmt.Println("Error:", err)
//		return
//	}
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
//	fmt.Println("SendReservationRequestApprovedBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
//	return
//}

func SendAdminReservationRequestApprovedBrevo(ctx context.Context, cfg *brevo.Configuration, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, expire string, templateID string, funcName string, appKey string) {
	toEmail := "sylviauwa@gmail.com"
	toName := "Sylvia"
	header := "Reservation request approved"
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

}

func SendOptionPaymentSuccessBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    guestFirstName,
		"name":         hostOptionName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"hostname":     hostName,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendOptionPaymentSuccessBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendOptionHostPaymentSuccessBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostFirstName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, guestEmail string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    hostFirstName,
		"name":         hostOptionName,
		"guestname":    guestFirstName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"guestemail":   guestEmail,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendOptionPaymentSuccessBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendReservationRequestApprovedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    guestFirstName,
		"name":         hostOptionName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"hostname":     hostName,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendOptionPaymentSuccessBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendReservationRequestDisapprovedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    guestFirstName,
		"name":         hostOptionName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"hostname":     hostName,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendOptionPaymentSuccessBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendDateUnavailableBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    guestFirstName,
		"name":         hostOptionName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"hostname":     hostName,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendDateUnavailableBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}

func SendPaymentFailedBrevo(ctx context.Context, cfg *brevo.Configuration, header string, hostOptionName string, hostName string, guestFirstName string, checkIn string, checkout string, adminPhoneNumber string, adminContactEmail string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"firstname":    guestFirstName,
		"name":         hostOptionName,
		"checkin":      checkIn,
		"checkout":     checkout,
		"hostname":     hostName,
		"phonenumber":  adminPhoneNumber,
		"contactemail": adminContactEmail,
		"year":         fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", header, message, hostOptionName)
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
	fmt.Println("SendPaymentFailedBrevo, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}
