package sender

import (
	"context"
	"fmt"
	"strconv"
	"time"

	brevo "github.com/getbrevo/brevo-go/lib"
)

func SendErrorMessageBrevo(ctx context.Context, cfg *brevo.Configuration, header string, message string, toEmail string, toName string, templateID string, funcName string, appKey string) (err error) {
	//cfg := brevo.NewConfiguration()
	////Configure API key authorization: api-key
	//cfg.AddDefaultHeader("api-key", appKey)
	params := map[string]any{
		"header":  header,
		"message": message,
		"time": fmt.Sprint(time.Now()),
		"year":    fmt.Sprint(time.Now().Year()),
	}
	dataString := fmt.Sprintf("<html><body><h1>%v</h1><div>%v</div><div></div><footer>Team Flizzup</footer></body></html", header, message)
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
