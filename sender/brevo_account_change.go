package sender

import (
	"context"
	"fmt"
	"strconv"
	"time"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendAccountChangeBrevo(ctx context.Context, cfg *brevo.Configuration, toName string, toEmail string, code string, templateID string, funcName string, appKey string, mainHeader string, header string, message string, expire string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"mainheader": mainHeader,
		"header": header,
		"message": message,
		"expire": expire,
		"code": code,
		"year": fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><h2>%v</h2><div>%v</div><div>%v</div><div>%v</div><footer>Team Flizzup</footer></body></html", mainHeader, header, message, code, expire)
	tID, err := strconv.Atoi(templateID)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error:", err)
		return
	}
	var paramsData any = params
	br := brevo.NewAPIClient(cfg)
	body := brevo.SendSmtpEmail{
		HtmlContent: dataString,
		TemplateId:  int64(tID),
		To: []brevo.SendSmtpEmailTo{
			{Email: toEmail, Name: toName},
		},
		ReplyTo: &brevo.SendSmtpEmailReplyTo{
			Name:  "team",
			Email: "support@flizzup.com",
		},
		Params: &paramsData,
	}
	obj, resp, err := br.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		fmt.Println("Error in TransactionalEmailsApi->SendTransacEmail ", err.Error(), "funcName: ", funcName)
		return
	}
	fmt.Println("SendTransacEmail, response:", resp, "SendTransacEmail object", obj, "funcName: ", funcName)
	return
}