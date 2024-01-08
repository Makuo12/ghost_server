package utils

import (
	"time"

	"github.com/spf13/viper"
)

//Config stores all the configuration of the application
//The values are read by viper from a config file or env variables
//to get the values of the variables and store then in the struct
//we need to use the unmarshal behavior of viper USING mapstructure
type Config struct {
	DBSource                      string        `mapstructure:"DB_Source"`
	MigrationURL                  string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress             string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress             string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey             string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	BrevoEmailTemplate            string        `mapstructure:"BREVO_EMAIL_TEMPLATE"`
	BrevoInviteTemplate           string        `mapstructure:"BREVO_INVITE_TEMPLATE"`
	BrevoCoHostDeactivateTemplate string        `mapstructure:"BREVO_COHOST_DEACTIVATE_TEMPLATE"`
	BrevoAccountChangeTemplate    string        `mapstructure:"BREVO_ACCOUNT_CHANGE_TEMPLATE"`
	BrevoReserveRequestTemplate   string        `mapstructure:"BREVO_RESERVE_REQUEST_TEMPLATE"`
	AccessTokenDuration           time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	BrevoApiKey                   string        `mapstructure:"BREVO_API_KEY"`
	RefreshTokenDuration          time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisPassword                 string        `mapstructure:"REDIS_PASSWORD"`
	RedisAddr                     string        `mapstructure:"REDIS_ADDR"`
	TwilioAccountSID              string        `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken               string        `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioPhoneNumber             string        `mapstructure:"TWILIO_PHONE_NUMBER"`
	SendGridAPIKey                string        `mapstructure:"SEND_GRID_API_KEY"`
	DollarToNaira                 string        `mapstructure:"DOLLAR_TO_NAIRA"`
	DollarToCAD                   string        `mapstructure:"DOLLAR_TO_CAD"`
	FirebaseAPIKeyFile            string        `mapstructure:"FIREBASE_API_KEY_FILE"`
	FirebaseBucketName            string        `mapstructure:"FIREBASE_BUCKET_NAME"`
	PaystackSecretTestKey         string        `mapstructure:"PAYSTACK_SECRET_TEST_KEY"`
	PaystackSecretLiveKey         string        `mapstructure:"PAYSTACK_SECRET_LIVE_KEY"`
	AddCardChargeNaira            string        `mapstructure:"ADD_CARD_CHARGE_NAIRA"`
	AddCardChargeDollar           string        `mapstructure:"ADD_CARD_CHARGE_DOLLAR"`
	IntServiceOptionUserPercent   string        `mapstructure:"INT_SERVICE_OPTION_USER_PERCENT"`
	LocServiceOptionUserPercent   string        `mapstructure:"LOC_SERVICE_OPTION_USER_PERCENT"`
	IntServiceOptionHostPercent   string        `mapstructure:"INT_SERVICE_OPTION_HOST_PERCENT"`
	LocServiceOptionHostPercent   string        `mapstructure:"LOC_SERVICE_OPTION_HOST_PERCENT"`
	IntServiceEventUserPercent    string        `mapstructure:"INT_SERVICE_EVENT_USER_PERCENT"`
	LocServiceEventUserPercent    string        `mapstructure:"LOC_SERVICE_EVENT_USER_PERCENT"`
	IntServiceEventHostPercent    string        `mapstructure:"INT_SERVICE_EVENT_HOST_PERCENT"`
	LocServiceEventHostPercent    string        `mapstructure:"LOC_SERVICE_EVENT_HOST_PERCENT"`
	PaymentSuccessUrl             string        `mapstructure:"PAYMENT_SUCCESS_URL"`
	PaymentFailUrl                string        `mapstructure:"PAYMENT_FAIL_URL"`
	WebsiteDomain                 string        `mapstructure:"WEBSITE_DOMAIN"`
	AppName                       string        `mapstructure:"APP_NAME"`
	AppEmail                      string        `mapstructure:"APP_EMAIL"`
	AppEmailName                  string        `mapstructure:"APP_EMAIL_NAME"`
	AppEmailDomain                string        `mapstructure:"APP_EMAIL_DOMAIN"`
	EmailTemplate                 string        `mapstructure:"EMAIL_TEMPLATE"`
	SmsOtpTemplate                string        `mapstructure:"SMS_OTP_TEMPLATE"`
	Msg91Key                      string        `mapstructure:"MSG91_KEY"`
	InviteTemplate                string        `mapstructure:"INVITE_TEMPLATE"`
	TermsOfService                string        `mapstructure:"TERMS_OF_SERVICE"`
	PrivacyPolicy                 string        `mapstructure:"PRIVACY_POLICY"`
}

func LoadConfig(path string) (config Config, err error) {
	//AddConfigPath tells viper the location of the env file
	viper.AddConfigPath(path)
	//SetConfigName would tell viper to search for a env file starting with app
	viper.SetConfigName("app")
	//SetConfigType allows us to tell viper the type the file because env,json,html can be used
	viper.SetConfigType("env")
	//AutomaticEnv would tell viper to automatically override
	//values that it has read from config file
	viper.AutomaticEnv()

	//ReadInConfig would start reading config values from config file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}
