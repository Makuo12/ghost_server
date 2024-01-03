package sender

type Variable struct {
	CompanyName string `json:"company_name"`
	Code        string `json:"code"`
	Year        string `json:"year"`
}

type To struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type From struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Recipient struct {
	To        []To     `json:"to"`
	Variables Variable `json:"variables"`
}

type EmailMessage struct {
	Recipients []Recipient `json:"recipients"`
	From       From        `json:"from"`
	Domain     string      `json:"domain"`
	TemplateID string      `json:"template_id"`
}

type VariableInvitation struct {
	CompanyName    string `json:"company_name"`
	Code           string `json:"code"`
	Year           string `json:"year"`
	Time           string `json:"time"`
	MainOption     string `json:"main_option"`
	MainHost       string `json:"main_host"`
	HostNameOption string `json:"host_name_option"`
}

type RecipientInvitation struct {
	To        []To               `json:"to"`
	Variables VariableInvitation `json:"variables"`
}

type EmailInvitationMessage struct {
	Recipients []RecipientInvitation `json:"recipients"`
	From       From                  `json:"from"`
	Domain     string                `json:"domain"`
	TemplateID string                `json:"template_id"`
}

type VariableCustom struct {
	CompanyName string `json:"company_name"`
	Year        string `json:"year"`
	Header      string `json:"header"`
	TopHeader   string `json:"top_header"`
	Body        string `json:"body"`
}

type RecipientCustom struct {
	To        []To           `json:"to"`
	Variables VariableCustom `json:"variables"`
}

type EmailCustomMessage struct {
	Recipients []RecipientCustom `json:"recipients"`
	From       From              `json:"from"`
	Domain     string            `json:"domain"`
	TemplateID string            `json:"template_id"`
}
