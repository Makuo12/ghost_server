package api

import (
	// "errors"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	//"log"
	// "net/http"
	db "flex_server/db/sqlc"
	"flex_server/token"
	"flex_server/utils"
	"flex_server/val"

	// "strings"
	//"time"

	//
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"

	brevo "github.com/getbrevo/brevo-go/lib"
	"google.golang.org/api/option"
)

//var Context *gin.Context
var RedisClient *redis.Client

type Server struct {
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
	store      *db.SQLStore
	Bucket     *storage.BucketHandle
	ClientFire *auth.Client
	ApnFire    *messaging.Client
	Cfg        *brevo.Configuration
}

var RedisContext = context.Background()

func NewServer(config utils.Config, store *db.SQLStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)

	}
	configFirebase := &firebase.Config{
		StorageBucket: config.FirebaseBucketName,
	}
	path, err := filepath.Abs(config.FirebaseAPIKeyFile)
	if err != nil {
		log.Printf("error you an error getting firebase json file %v\n", err.Error())
	}
	opt := option.WithCredentialsFile(path)
	appFirebase, err := firebase.NewApp(context.Background(), configFirebase, opt)
	if err != nil {
		log.Fatalln(err)
	}

	clientFirebase, err := appFirebase.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	clientFire, err := appFirebase.Auth(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	apnFire, err := appFirebase.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	bucket, err := clientFirebase.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}
	cfg := brevo.NewConfiguration()
	//Configure API key authorization: api-key
	cfg.AddDefaultHeader("api-key", config.BrevoApiKey)
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		Bucket:     bucket,
		ClientFire: clientFire,
		ApnFire:    apnFire,
		Cfg:        cfg,
	}
	// adr := []string{"redis_flex:6379"}
	// rdb := redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName: "Flex",
	// 	SentinelAddrs: adr,
	// 	SentinelPassword: "Si73gangan",
	// 	DB: 0,
	// })
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis_flex:6379",
		Password: server.config.RedisPassword, // no password set
		DB:       0,                           // use default DB
	})
	RedisClient = rdb
	rdb.Ping(RedisContext)

	//RedisClient.FlushAll(RedisContext)

	// We setup up the validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		err = v.RegisterValidation("currency", validCurrency)
		if err != nil {
			log.Printf("Error at server setup in, validCurrency Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("email", validEmail)
		if err != nil {
			log.Printf("Error at server setup, validEmail in Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("password", validPassword)
		if err != nil {
			log.Printf("Error at server setup, validPassword in Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("time_only", validTimeOnly)
		if err != nil {
			log.Printf("Error at server setup, validTimeOnly in Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("date_only", validDateOnly)
		if err != nil {
			log.Printf("Error at server setup, validDateOnly in Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("money", validMoney)
		if err != nil {
			log.Printf("Error at server setup, validMoney in Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("person_name", validName)
		if err != nil {
			log.Printf("Error at server setup in, validName Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("shortlet_space", validShortletSpace)
		if err != nil {
			log.Printf("Error at server setup in, validShortletSpace Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("shortlet_type", validShortletType)
		if err != nil {
			log.Printf("Error at server setup in, validShortletType Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("event_location_type", validEventLocationType)
		if err != nil {
			log.Printf("Error at server setup in, validEventLocationType Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("event_date_type", validEventDateType)
		if err != nil {
			log.Printf("Error at server setup in, validEventDateType Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("event_ticket_type", validEventTicketType)
		if err != nil {
			log.Printf("Error at server setup in, validEventTicketType Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("event_ticket_level", validEventTicketLevel)
		if err != nil {
			log.Printf("Error at server setup in, validEventTicketLevel Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("event_ticket_main_type", validEventTicketMainType)
		if err != nil {
			log.Printf("Error at server setup in, validEventTicketMainType Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("property_size_unit", validPropertySizeUnit)
		if err != nil {
			log.Printf("Error at server setup in, validPropertySizeUnit Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("des_type", validDesTypes)
		if err != nil {
			log.Printf("Error at server setup in, validDesTypes Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("option_extra_info_type", validOptionExtraInfo)
		if err != nil {
			log.Printf("Error at server setup in, validOptionExtraInfo Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("user_profile_type", validUserProfileTypes)
		if err != nil {
			log.Printf("Error at server setup in, validUserProfileTypes Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("report_option", validReportOption)
		if err != nil {
			log.Printf("Error at server setup in,  validReportOption Register Validation %v", err.Error())
		}
		err = v.RegisterValidation("guest_option", validGuestOption)
		if err != nil {
			log.Printf("Error at server setup in, validGuestOption Register Validation %v", err.Error())
		}

		err = v.RegisterValidation("user_option_cancel", validUserOptionCancel)
		if err != nil {
			log.Printf("Error at server setup in, validUserOptionCancel Register Validation %v", err.Error())
		}

		err = v.RegisterValidation("user_event_cancel", validUserEventCancel)
		if err != nil {
			log.Printf("Error at server setup in, validUserEventCancel Register Validation %v", err.Error())
		}

		err = v.RegisterValidation("host_option_cancel", validHostOptionCancel)
		if err != nil {
			log.Printf("Error at server setup in, validHostOptionCancel Register Validation %v", err.Error())
		}

		err = v.RegisterValidation("host_event_cancel", validHostEventCancel)
		if err != nil {
			log.Printf("Error at server setup in, validHostEventCancel Register Validation %v", err.Error())
		}
	}
	ctx := context.Background()

	job := cron.New()

	// Schedule the daily function to run at 5 hours
	_, err = job.AddFunc("@every 2m", DailyRemoveOptionReserveUser)
	if err != nil {
		log.Printf("Error at cron at job.AddFunc for DailyRemoveOptionReserveUser %v", err.Error())
	}
	// Schedule the daily function to run at 5 hours
	_, err = job.AddFunc("@every 2m", DailyRemoveEventReserveUser)
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyRemoveEventReserveUser  %v", err.Error())
	}

	// DailyHandleRefund
	_, err = job.AddFunc("@every 2m", DailyHandleRefund(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandleTransferWebhookData  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyHandleTransferWebhookData(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandleTransferWebhookData  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyHandleRefundWebhookData(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for  DailyHandleRefundWebhookData  %v", err.Error())
	}

	// Payout
	_, err = job.AddFunc("@every 2m", DailyHandlePayouts(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandlePayouts  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyHandleRefundPayouts(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandleRefundPayouts  %v", err.Error())
	}

	// Event date change
	_, err = job.AddFunc("@every 1m", DailyChangeDateEventHostUpdate(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyChangeDateEventHostUpdate  %v", err.Error())
	}

	// Event date cancellation
	_, err = job.AddFunc("@every 1m", DailyCreateEventHostCancel(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyCreateEventHostCancel  %v", err.Error())
	}

	//// Snooze
	_, err = job.AddFunc("@every 2m", DailyHandleSnooze(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandleSnooze  %v", err.Error())
	}

	//// User Request
	_, err = job.AddFunc("@every 2m", DailyHandleUserRequest(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyHandleUserRequest  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyDeactivateCoHost(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyDeactivateCoHost  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyValidatedChargeTicket(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyValidatedChargeTicket  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", DailyValidatedChargeOption(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for DailyValidatedChargeOption  %v", err.Error())
	}

	//// Experience
	//HandleOptionExperienceToRedis(ctx, server)
	//HandleEventExperienceToRedis(ctx, server)

	_, err = job.AddFunc("@every 2m", HandleOptionExperienceToRedis(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for HandleOptionExperienceToRedis  %v", err.Error())
	}

	_, err = job.AddFunc("@every 2m", HandleEventExperienceToRedis(ctx, server))
	if err != nil {
		log.Printf(" Error at cron at job.AddFunc for HandleEventExperienceToRedis  %v", err.Error())
	}

	//_, err = job.AddFunc("@every 2m", HandleImageMetaData(ctx, server))
	//if err != nil {
	//	log.Printf(" Error at cron at job.AddFunc for HandleImageMetaData  %v", err.Error())
	//}

	// Start the cron scheduler
	job.Start()

	server.setupRouter()

	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())

	//err := router.SetTrustedProxies([]string{"172.18.0.1", "192.168.48.1", "52.31.139.75", "52.49.173.169", "52.214.14.220"})
	//if err != nil {
	//	log.Printf("Error at server setupRouter %v", err.Error())
	//}

	router.POST("/users/email/check", server.JoinVerifyEmail)
	router.POST("/users/email/confirm", server.ConfirmCode)

	router.POST("/users/phone/login/check", server.ForgotPasswordNotLogged)
	router.POST("/users/phone/login/confirm", server.ConfirmCodeLogin)
	router.POST("/users/create", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	router.POST("/tokens/renew_access", server.RenewAccessToken)

	router.PUT("/users/get-app-policy", server.GetAppPolicy)

	router.PUT("/users/join/phone", server.JoinWithPhone)
	router.PUT("/users/join/sign-up/confirm", server.ConfirmJoinSignUp)
	router.POST("/users/join/code/confirm", server.ConfirmCodeJoin)

	// Experience
	router.PUT("/users/user/experience/list", server.ListExperience)
	router.PUT("/users/user/experience/detail", server.GetExperienceDetail)
	router.PUT("/users/user/experience/detail/tickets", server.ListExperienceEventTickets)
	router.PUT("/users/user/experience/option/am-detail/get", server.GetExperienceAmDetail)
	router.PUT("/users/user/experience/detail/option/date-time/list", server.ListExOptionDateTime)

	// Experience Search
	router.PUT("/users/user/experience/options/search-filter", server.ListOptionExSearch)
	router.PUT("/users/user/experience/events/search-filter", server.ListEventExSearch)
	// Filter Range
	router.PUT("/users/user/experience/options/filter-range", server.GetOptionFilterRange)
	router.PUT("/users/user/experience/events/filter-range", server.GetEventFilterRange)

	// Deep Link
	router.PUT("/users/user/experience/deep-link/option/get", server.GetOptionDeepLinkExperience)
	router.PUT("/users/user/experience/deep-link/event/get", server.GetEventDeepLinkExperience)
	router.PUT("/users/user/experience/deep-link/event/event-dates/get", server.GetEventDateDeepLinkExperience)

	// Update without being logged in
	// lock means password
	router.PUT("/users/not_logged/forgot/lock", server.ForgotPasswordNotLogged)
	router.PUT("/users/not_logged/forgot/confirm-code", server.ConfirmCode)
	router.PUT("/users/not_logged/forgot/new-lock", server.NewPassword)
	router.POST("/users/not-user/support/help/create", server.CreateHelp)

	// APN TOKEN
	router.POST("/users/user/apn/create-detail", server.CreateUserAPNDetail)

	// Webhook
	router.POST("/webhook/paystack", server.PaystackWebhook)
	// WEBSOCKET
	router.GET("/ws/user/:app/:room_type", func(ctx *gin.Context) {
		roomType := ctx.Param("room_type")
		app := ctx.Param("app")
		result := val.ContainRoomType(roomType)

		if !result {
			err := fmt.Errorf("this path does not exist")
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
		}
		if roomType != "map_location" && roomType != "ex_search_event" {
			err := fmt.Errorf("this path does not exist")
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
		}
		ServeWs(ctx.Writer, ctx.Request, roomType, uuid.New(), uuid.New(), app, server, ctx, db.User{})
	})

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// Create options
	authRoutes.POST("/users/options_se/create", server.CreateOptionSE)
	authRoutes.DELETE("/users/options_se/remove", server.RemoveOptionSE)

	authRoutes.POST("/users/location/create", server.CreateLocation)
	authRoutes.DELETE("/users/location/remove", server.RemoveLocation)

	authRoutes.POST("/users/event-sub-type/create", server.CreateEventSubCategory)
	authRoutes.DELETE("/users/event-sub-type/remove", server.RemoveEventSubCategory)

	authRoutes.POST("/users/shortlet-space/create", server.CreateShortletSpace)
	authRoutes.PUT("/users/shortlet-space/remove", server.RemoveShortletSpace)

	authRoutes.POST("/users/shortlet-info/create", server.CreateShortletInfo)
	authRoutes.PUT("/users/shortlet-info/remove", server.RemoveShortletInfo)

	authRoutes.POST("/users/amenities/create", server.CreateAmenitiesAndSafety)
	authRoutes.DELETE("/users/amenities/remove", server.RemoveAmenitiesAndSafety)

	authRoutes.POST("/users/description/create", server.CreateOptionDescription)
	authRoutes.DELETE("/users/description/remove", server.RemoveOptionDescription)

	authRoutes.POST("/users/option-name/create", server.CreateOptionName)
	authRoutes.DELETE("/users/option-name/remove", server.RemoveOptionName)

	authRoutes.POST("/users/option-price/create", server.CreateOptionPrice)
	authRoutes.DELETE("/users/option-price/remove", server.RemoveOptionPrice)

	authRoutes.POST("/users/option-highlight/create", server.CreateOptionHighlight)
	authRoutes.DELETE("/users/option-highlight/remove", server.RemoveOptionHighlight)

	authRoutes.POST("/users/option-photo/create", server.CreateOptionPhoto)
	authRoutes.DELETE("/users/option-photo/remove", server.RemoveOptionPhoto)

	authRoutes.POST("/users/option-question/create", server.CreateOptionQuestion)
	authRoutes.DELETE("/users/option-question/remove", server.RemoveOptionQuestion)

	authRoutes.POST("/users/publish-option/create", server.PublishOption)
	authRoutes.DELETE("/users/publish-option/remove", server.RemovePublishOption)
	authRoutes.PUT("/users/publish-option/get", server.GetPublishData)
	authRoutes.DELETE("/users/logout-user", server.LogoutUser)

	// End of Create option or event

	// Start of calenderView routes
	// MAIN OPTIONS: REFER TO SHORTLET OR EVENTS
	authRoutes.PUT("/users/calender-options/list", server.ListCalenderOptionItems)

	// EVENT DATE TIME
	authRoutes.PUT("/users/calender-options/event/date-time/list", server.ListEventDateItems)
	authRoutes.PUT("/users/calender-options/event/date-time/normal/list", server.ListEventDateNormalItems)

	authRoutes.POST("/users/calender-options/event/date-time/create", server.CreateEventDateTime)
	authRoutes.PUT("/users/calender-options/event/date-time/update", server.UpdateEventDateTime)
	authRoutes.PUT("/users/calender-options/event/date-time/controls/update", server.UpdateEventDateTimeControls)
	authRoutes.PUT("/users/calender-options/event/date-time/note/update", server.UpdateEventDateTimeNote)
	authRoutes.DELETE("/users/calender-options/event/date-time/remove", server.RemoveEventDateTime)

	// EVENT DATE STATUS
	authRoutes.PUT("/users/calender-options/event/status/update", server.UpdateEventDateStatus)
	authRoutes.PUT("/users/calender-options/event/status/get", server.GetEventDateIsBooked)
	authRoutes.PUT("/users/calender-options/event/dates-booked/list", server.ListEventDateBooked)
	authRoutes.PUT("/users/reservation/event/reserve/change-dates/update", server.UpdateEventDatesBooked)

	// EVENT DATE TICKET
	authRoutes.POST("/users/calender-options/event/ticket/create", server.CreateEventDateTicket)
	authRoutes.PUT("/users/calender-options/event/ticket/update", server.UpdateEventDateTicket)
	authRoutes.PUT("/users/calender-options/event/ticket/list", server.ListEventDateTicket)
	authRoutes.DELETE("/users/calender-options/event/ticket/remove", server.RemoveEventDateTicket)

	// EVENT DATE LOCATION
	authRoutes.POST("/users/calender-options/event/location/create-update", server.CreateUpdateEventDateLocation)
	authRoutes.PUT("/users/calender-options/event/location/get", server.GetEventDateLocation)

	// EVENT DATE DETAIL
	authRoutes.POST("/users/calender-options/event/detail/create-update", server.CreateUpdateEventDateDetail)

	// EVENT DATE PUBLISHES
	authRoutes.PUT("/users/calender-options/event/publish/get", server.GetEventDatePublish)
	authRoutes.PUT("/users/calender-options/event/publish/update", server.UpdateEventDatePublish)

	// Audiences
	authRoutes.PUT("/users/calender-options/event/publish/audience/update", server.UpdatePrivateAudience)
	authRoutes.DELETE("/users/calender-options/event/publish/audience/remove", server.RemovePrivateAudience)

	// OPTIONS ONLY
	authRoutes.PUT("/users/calender-options/option/date-time/list", server.ListOptionDateItems)

	authRoutes.POST("/users/calender-options/option/date-time/create-update", server.CreateUpdateOptionDateTime)

	authRoutes.PUT("/users/calender-options/option/date-time/get-note", server.GetOptionDateNote)
	authRoutes.DELETE("/users/calender-options/option/date-time/remove", server.RemoveOptionDateTime)

	// End of calender routes

	authRoutes.PUT("/users/logged/forgot/lock", server.ForgotPasswordLogged)
	authRoutes.PUT("/users/logged/forgot/confirm-code", server.ConfirmCode)
	authRoutes.PUT("/users/logged/forgot/new-lock", server.NewPassword)
	//

	// User manage option selection
	authRoutes.PUT("/users/uhm/selection/get", server.ListUHMOptionSelection)
	// End User manage option selection

	// Start of User host manage options routes

	// SPACE AREAS
	authRoutes.PUT("/users/uhm/uhm-options/get", server.GetUHMOptionData)
	authRoutes.PUT("/users/uhm/space-areas/list", server.ListSpaceArea)
	authRoutes.PUT("/users/uhm/space-areas/create-edit", server.CreateEditSpaceAreas)
	authRoutes.PUT("/users/uhm/space-areas/add-bed/update", server.AddBedSpaceAreas)
	authRoutes.PUT("/users/uhm/space-areas/add-photo/update", server.AddPhotoSpaceAreas)
	authRoutes.PUT("/users/uhm/space-areas/update", server.UpdateSpaceAreas)

	// SHORTLET INFO
	authRoutes.PUT("/users/uhm/shortlet-info/update", server.UpdateShortletInfo)
	authRoutes.PUT("/users/uhm/shortlet-info/get", server.GetShortletInfo)

	// EVENT INFO
	authRoutes.PUT("/users/uhm/event-info/update", server.UpdateEventInfo)
	authRoutes.PUT("/users/uhm/event-info/get", server.GetEventInfo)

	// TITLE
	authRoutes.PUT("/users/uhm/option-detail/title/update", server.UpdateOptionTitle)

	authRoutes.PUT("/users/uhm/option-detail/des/update", server.UpdateOptionDes)
	authRoutes.PUT("/users/uhm/option-detail/des/get", server.GetOptionDes)

	// PHOTO
	authRoutes.PUT("/users/uhm/option-photo/photo/:photo_type", server.UpdateOptionPhoto)
	authRoutes.PUT("/users/uhm/option-photo-caption/create-update", server.CreateUpdateOptionPhotoCaption)
	authRoutes.PUT("/users/uhm/option-photo-caption/get", server.GetOptionPhotoCaption)

	//LOCATION
	authRoutes.PUT("/users/uhm/option-location/get", server.GetOptionLocation)
	authRoutes.PUT("/users/uhm/option-location/update", server.UpdateOptionLocation)
	authRoutes.PUT("/users/uhm/option-location/specific-location/update", server.UpdateShowSpecificLocation)

	//Highlight
	authRoutes.PUT("/users/uhm/option-detail/highlight/get", server.GetUHMHighlight)
	authRoutes.PUT("/users/uhm/option-detail/highlight/update", server.UpdateUHMHighlight)

	// EXTRA INFO
	authRoutes.PUT("/users/uhm/option-extra-info/get", server.GetOptionExtraInfo)
	authRoutes.PUT("/users/uhm/option-extra-info/create-update", server.CreateUpdateOptionExtraInfo)

	// WIFI
	authRoutes.PUT("/users/uhm/wifi-detail/get", server.GetOptionWifiDetail)
	authRoutes.PUT("/users/uhm/wifi-detail/create-update", server.CreateUpdateWifiDetail)

	// AMENITIES
	authRoutes.PUT("/users/uhm/amenities/list", server.ListUHMAmenities)
	authRoutes.PUT("/users/uhm/amenities/create-update", server.CreateUpdateAmenity)
	authRoutes.PUT("/users/uhm/amenities/detail/get", server.GetAmenityDetail)
	authRoutes.PUT("/users/uhm/amenities/detail/update", server.UpdateAmenityDetail)

	// CHECK IN STEP
	// Shortlet
	authRoutes.PUT("/users/uhm/check-in-step/list", server.ListCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/photo/remove", server.RemoveCheckInStepPhoto)
	authRoutes.PUT("/users/uhm/check-in-step/update", server.UpdateCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/create", server.CreateCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/remove", server.RemoveCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/publish/update", server.UpdatePublishOptionCheckInStep)
	// Event
	authRoutes.PUT("/users/uhm/check-in-step/event/list", server.ListEventCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/event/photo/remove", server.RemoveEventCheckInStepPhoto)
	authRoutes.PUT("/users/uhm/check-in-step/event/update", server.UpdateEventCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/event/create", server.CreateEventCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/event/remove", server.RemoveEventCheckInStep)
	authRoutes.PUT("/users/uhm/check-in-step/event/publish/update", server.UpdatePublishEventCheckInStep)

	// Additional Charge
	authRoutes.PUT("/users/uhm/booking/option-add-charge/list", server.ListOptionAddCharge)
	authRoutes.PUT("/users/uhm/booking/option-add-charge/create-update", server.CreateUpdateOptionAddCharge)

	// Pets Allowed
	authRoutes.PUT("/users/uhm/pets-allowed/update", server.UpdatePetsAllowed)

	// Length Of Stay Discount
	authRoutes.PUT("/users/uhm/booking/length-of-stay/option-discount/list", server.LOTListOptionDiscount)
	authRoutes.PUT("/users/uhm/booking/length-of-stay/option-discount/create-update", server.LOTCreateUpdateOptionDiscount)

	// Price
	authRoutes.PUT("/users/uhm/booking/option-price/update", server.UpdateOptionPrice)
	authRoutes.PUT("/users/uhm/booking/option-price/get", server.GetOptionPrice)

	// Currency
	authRoutes.PUT("/users/uhm/booking/currency/update", server.UpdateOptionCurrency)

	// Currency for unlisted options
	authRoutes.PUT("/users/uhm/booking/unlisted/currency/update", server.UpdateUnlistedOptionCurrency)

	// CHECK IN METHOD
	authRoutes.PUT("/users/uhm/shortlet/check-in-method/get", server.GetShortletCheckInMethod)
	authRoutes.PUT("/users/uhm/shortlet/check-in-method/update", server.UpdateShortletCheckInMethod)

	// AVAILABILITY SETTINGS
	authRoutes.PUT("/users/uhm/shortlet/available_setting/get", server.GetOptionAvailabilitySetting)
	authRoutes.PUT("/users/uhm/shortlet/available_setting/update", server.UOptionAvailabilitySetting)

	// TRIP LENGTH
	authRoutes.PUT("/users/uhm/shortlet/trip_length/get", server.GetOptionTripLength)
	authRoutes.PUT("/users/uhm/shortlet/trip_length/update", server.UOptionTripLength)

	// CHECK IN OUT DETAIL
	authRoutes.PUT("/users/uhm/shortlet/check_in_out_detail/get", server.GetCheckInOutDetail)
	authRoutes.PUT("/users/uhm/shortlet/check_in_out_detail/update", server.UCheckInOutDetail)

	// CANCEL POLICY
	authRoutes.PUT("/users/uhm/shortlet/cancel_policy/get", server.GetCancelPolicy)
	authRoutes.PUT("/users/uhm/shortlet/cancel_policy/update", server.UpdateCancelPolicy)

	// BOOK REQUIREMENT
	authRoutes.PUT("/users/uhm/shortlet/book_requirement/get", server.GetBookRequirement)
	authRoutes.PUT("/users/uhm/shortlet/book_requirement/update", server.UBookRequirement)

	// BOOK METHOD
	authRoutes.PUT("/users/uhm/shortlet/book_method/get", server.GetOptionBookMethod)
	authRoutes.PUT("/users/uhm/shortlet/book_method/update", server.UOptionBookMethod)
	authRoutes.PUT("/users/uhm/shortlet/book_method/msg/update", server.UOptionBookMethodMsg)

	// OPTION STATUS
	authRoutes.PUT("/users/uhm/shortlet/option_status/get", server.GetOptionInfoStatus)
	authRoutes.PUT("/users/uhm/shortlet/option_status/update", server.UpdateOptionInfoStatus)

	// THINGS TO NOTE
	authRoutes.PUT("/users/uhm/shortlet/thing-to-note/create-update", server.CUThingToNote)
	authRoutes.PUT("/users/uhm/shortlet/thing-to-note/list", server.ListThingToNote)
	authRoutes.PUT("/users/uhm/shortlet/thing-to-note/detail/get", server.GetThingToNoteDetail)
	authRoutes.PUT("/users/uhm/shortlet/thing-to-note/detail/update", server.UThingToNoteDetail)

	// House rules
	authRoutes.PUT("/users/uhm/option-rule/create-update", server.CUOptionRule)
	authRoutes.PUT("/users/uhm/option-rule/list", server.ListOptionRule)
	authRoutes.PUT("/users/uhm/option-rule/detail/get", server.GetOptionRuleDetail)
	authRoutes.PUT("/users/uhm/option-rule/detail/update", server.UOptionRuleDetail)

	// Co-Host Primary
	authRoutes.POST("/users/uhm/option-co-host/create", server.CreateOptionCoHost)
	authRoutes.PUT("/users/uhm/option-co-host/update", server.UpdateOptionCoHost)
	authRoutes.PUT("/users/uhm/option-co-host/get", server.GetOptionCoHost)
	authRoutes.PUT("/users/uhm/option-co-host/list", server.ListOptionCoHost)
	authRoutes.PUT("/users/uhm/option-co-host/invite/resend", server.ResendInviteCoHost)
	authRoutes.DELETE("/users/uhm/option-co-host/invite/delete", server.CancelInviteCoHost)
	authRoutes.DELETE("/users/uhm/option-co-host/invite/remove", server.RemoveCoHost)
	// End of user manage option routes
	// Get User
	authRoutes.GET("/users/fire-fight/get", server.GetFireEmailAndPassword)

	// Update Users
	authRoutes.PUT("/users/currency/update", server.UpdateCurrency)
	authRoutes.PUT("/users/user/update/info", server.UpdateUserInfo)
	authRoutes.PUT("/users/user/update/password", server.UpdateUserPassword)
	authRoutes.PUT("/users/user/update/verify", server.UpdateVerifyEmailPhone)
	authRoutes.PUT("/users/user/update/currency", server.UpdateCurrencyUser)
	authRoutes.PUT("/users/user/update/code", server.UpdateCodeEmailPhone)
	authRoutes.GET("/users/user/get/profile-user", server.GetProfileUser)
	authRoutes.DELETE("/users/user/profile-user/em-contact/remove", server.RemoveEmContact)
	authRoutes.POST("/users/user/profile-user/em-contact/create", server.CreateEmContact)
	authRoutes.POST("/users/user/profile-user/fire/identity/update", server.UpdateIdentity)
	authRoutes.GET("/users/user/currency/get", server.GetUserCurrency)
	authRoutes.GET("/users/user/profile-photo/get", server.GetUserProfilePhoto)
	authRoutes.PUT("/users/user/profile-photo/", server.GetUserCurrency)

	// Profile Detail
	authRoutes.GET("/users/user/get/profile-detail/get", server.GetUserProfileDetail)
	authRoutes.PUT("/users/user/profile-detail/user-profile/update", server.UpdateUserProfile)
	authRoutes.PUT("/users/user/profile-detail/profile-photo/update", server.UpdateUserProfilePhoto)
	authRoutes.PUT("/users/user/profile-detail/user-location/create-update", server.CreateUpdateUserLocation)
	authRoutes.DELETE("/users/user/profile-detail/user-location/remove", server.RemoveUserLocation)

	// Support
	authRoutes.POST("/users/user/support/feedback/create", server.CreateFeedback)
	authRoutes.POST("/users/user/support/help/create", server.CreateHelpUser)

	// Report Option
	authRoutes.POST("/users/user/report/option/create", server.CreateReportOptionUser)

	// User
	authRoutes.GET("/users/user/is-host/get", server.GetUserIsHost)
	authRoutes.GET("/users/user/get-start", server.GetUser)

	// get user options that are incomplete
	authRoutes.PUT("/users/user/option-info/not-complete/list", server.ListIncompleteOptionInfos)
	authRoutes.GET("/users/user/option-info/not-complete/get/:option_id", server.GetOptionInfoIncomplete)
	//

	// Payments

	authRoutes.POST("/users/init/add/reference", server.VerifyAddCardChargeReference)
	authRoutes.POST("/users/user/init/add", server.InitAddCard)
	authRoutes.DELETE("/users/user/init/remove", server.InitRemoveCard)
	authRoutes.PUT("/users/user/init/default/update", server.SetDefaultCard)
	authRoutes.GET("/users/user/init/wall/get", server.GetWallet)

	// Direct Payments
	authRoutes.POST("/users/init/direct/payment", server.InitPayment)
	authRoutes.POST("/users/init/direct/payment/reference", server.VerifyPaymentReference)

	// Reservation
	////-> Option
	authRoutes.POST("/users/user/reserve/option", server.CreateOptionReserveDetail)
	authRoutes.POST("/users/user/reserve/option/final/create", server.FinalOptionReserveDetail)
	authRoutes.POST("/users/user/reserve/option/final-verification/create", server.FinalOptionReserveVerificationDetail)

	////-> Event
	authRoutes.POST("/users/user/reserve/event", server.CreateEventReserveDetail)
	authRoutes.POST("/users/user/reserve/event/final/create", server.FinalEventReserveDetail)
	authRoutes.POST("/users/user/reserve/event/final-verification/create", server.FinalEventReserveVerificationDetail)

	// Reservation Host Section
	authRoutes.GET("/users/host/reserve/get", server.GetReserveHostDetail)
	authRoutes.PUT("/users/host/reserve/list", server.ListReservationDetail)

	// Messages
	authRoutes.PUT("/users/user/message/contact/list", server.ListMessageContact)
	authRoutes.PUT("/users/user/message/message/list", server.ListMessage)
	authRoutes.POST("/users/user/message/message/create", server.CreateMessage)

	// Notifications
	authRoutes.PUT("/users/user/notification/list", server.ListNotification)
	authRoutes.PUT("/users/user/notification/detail", server.NotificationOptionReserveDetail)

	// Request Notify
	authRoutes.PUT("/users/user/request-notify/list", server.ListRequestNotify)
	authRoutes.PUT("/users/user/option/request-notify/user-request/get", server.GetOUserRequestNotifyDetail)
	authRoutes.PUT("/users/user/option/request-notify/ans-approve", server.MsgRequestResponse)

	// Account Number
	authRoutes.PUT("/users/user/bank/list", server.ListBank)
	authRoutes.GET("/users/user/account-number/list", server.ListAccountNumber)
	authRoutes.POST("/users/user/account-number/create", server.CreateAccountNumber)
	authRoutes.DELETE("/users/user/account-number/remove", server.RemoveAccountNumber)
	authRoutes.PUT("/users/user/account-number/set-default", server.SetDefaultAccountNumber)

	// Wishlist
	authRoutes.POST("/users/user/wishlist/create", server.CreateWishlist)
	authRoutes.POST("/users/user/wishlist/item/create", server.CreateWishlistItem)
	authRoutes.DELETE("/users/user/wishlist/remove", server.RemoveWishlist)
	authRoutes.DELETE("/users/user/wishlist/item/remove", server.RemoveWishlistItem)
	authRoutes.GET("/users/user/wishlist/list", server.ListWishlist)
	authRoutes.PUT("/users/user/experience/wishlist/list", server.ListWishlistExperience)

	// RESERVE USER
	authRoutes.PUT("/users/user/reserve-user/list", server.ListReserveUserItem)
	authRoutes.PUT("/users/user/reserve-user/direction", server.GetReserveUserDirection)
	authRoutes.PUT("/users/user/reserve-user/check-in-step", server.GetRUCheckInStep)
	authRoutes.PUT("/users/user/reserve-user/help", server.GetRUHelp)
	authRoutes.PUT("/users/user/reserve-user/wifi", server.GetRUWifi)
	authRoutes.PUT("/users/user/reserve-user/receipt", server.GetRUReceipt)

	// PAYOUT AND REFUNDS
	authRoutes.PUT("/users/user/payment/event/payout/list", server.ListEventPayout)
	authRoutes.PUT("/users/user/payment/option/payout/list", server.ListOptionPayout)
	authRoutes.PUT("/users/user/payment/option/payment/list", server.ListOptionPayment)
	authRoutes.PUT("/users/user/payment/ticket/payment/list", server.ListTicketPayment)
	authRoutes.PUT("/users/user/payment/refund/list", server.ListRefund)
	authRoutes.PUT("/users/user/payment/refund-payout/list", server.ListRefundPayout)

	// Co-Host Secondary
	authRoutes.PUT("/users/user/option-co-host/co-host-user/list", server.ListOptionCoHostItem)
	authRoutes.PUT("/users/user/option-co-host/co-host-user/invite/validate", server.ValidateOptionCoHost)
	authRoutes.PUT("/users/user/option-co-host/co-host-user/detail", server.GetOptionCoHostItemDetail)
	authRoutes.DELETE("/users/user/option-co-host/co-host-user/delete", server.DeactivateOptionCoHost)

	// Cancellation User
	authRoutes.POST("/users/reserve-user/option/user-cancel/create", server.CreateOptionUserCancel)
	authRoutes.POST("/users/reserve-user/event/user-cancel/create", server.CreateEventUserCancel)

	// Cancellation Host
	authRoutes.POST("/users/reserve/option/host-cancel/create", server.CreateOptionHostCancel)
	authRoutes.POST("/users/reserve/event/host-cancel/create", server.CreateEventHostCancel)

	// Cancellation User Detail
	authRoutes.PUT("/users/reserve-user/option/user-cancel/detail/get", server.GetOptionUserCancelDetail)
	authRoutes.PUT("/users/reserve-user/event/user-cancel/detail/get", server.GetEventUserCancelDetail)

	// Cancellation Host Detail
	authRoutes.PUT("/users/reserve/option/host-cancel/detail/get", server.GetOptionHostCancelDetail)
	authRoutes.PUT("/users/reserve/event/host-cancel/detail/get", server.GetEventHostCancelDetail)

	// Scan Code
	authRoutes.PUT("/users/user/reserve/scan-code/charge/get", server.GetChargeCode)
	authRoutes.DELETE("/users/user/reserve/scan-code/charge/delete", server.DeleteChargeCode)
	authRoutes.PUT("/users/user/host/reserve/scan-code/charge/validate", server.ValidateChargeCode)

	// Insights
	authRoutes.PUT("/users/option/insight/list", server.ListOptionInsight)
	// Option Insights
	authRoutes.PUT("/users/option/insight/shortlet/get", server.GetOptionInsight)
	authRoutes.PUT("/users/option/insight/shortlet/all/get", server.GetAllOptionInsight)

	// Event Insights
	authRoutes.PUT("/users/option/insight/event/get", server.GetEventInsight)
	authRoutes.PUT("/users/option/insight/event/all/get", server.GetAllEventInsight)

	// Reviews
	authRoutes.PUT("/users/reviews/user/get-state", server.GetStateReview)
	authRoutes.POST("/users/reviews/user/create-general", server.CreateGeneralReview)
	authRoutes.PUT("/users/reviews/user/create-detail", server.CreateDetailReview)
	authRoutes.PUT("/users/reviews/user/create-private-note", server.CreatePrivateNoteReview)
	authRoutes.PUT("/users/reviews/user/create-public-note", server.CreatePublicNoteReview)
	authRoutes.PUT("/users/reviews/user/extra-placeholder", server.PlaceholderReview)
	authRoutes.PUT("/users/reviews/user/create-stay-clean", server.CreateStayCleanReview)
	authRoutes.PUT("/users/reviews/user/create-comfort", server.CreateComfortReview)
	authRoutes.PUT("/users/reviews/user/create-host", server.CreateHostReview)
	authRoutes.PUT("/users/reviews/user/list-amenity", server.ListAmenityReview)
	authRoutes.PUT("/users/reviews/user/create-amenity", server.CreateAmenityReview)
	authRoutes.PUT("/users/reviews/user/remove-amenity", server.RemoveAmenityReview)
	authRoutes.PUT("/users/reviews/user/create-option", server.CompleteOptionReview)

	// Remove Reviews
	authRoutes.DELETE("/users/reviews/user/remove-general", server.RemoveGeneralReview)
	authRoutes.PUT("/users/reviews/user/remove-detail", server.RemoveDetailReview)
	authRoutes.PUT("/users/reviews/user/remove-private-note", server.RemovePrivateNoteReview)
	authRoutes.PUT("/users/reviews/user/remove-public-note", server.RemovePublicNoteReview)
	authRoutes.PUT("/users/reviews/user/remove-stay-clean", server.RemoveStayCleanReview)
	authRoutes.PUT("/users/reviews/user/remove-comfort", server.RemoveComfortReview)
	authRoutes.PUT("/users/reviews/user/remove-host", server.RemoveHostReview)
	authRoutes.PUT("/users/reviews/user/remove-all-amenity", server.RemoveAllAmenityReview)

	authRoutes.PUT("/users/account-change/user/change", server.AccountChange)
	authRoutes.PUT("/users/account-change/user/verify-change", server.VerifyAccountChange)

	// Deep links
	authRoutes.PUT("/users/options/deep-link/get", server.GetOptionDeepLink)
	authRoutes.PUT("/users/options/event-dates/deep-link/get", server.GetEventDateDeepLink)

	// Option Questions
	authRoutes.PUT("/users/options/option-question/get", server.GetOptionQuestion)
	authRoutes.PUT("/users/options/option-question/update", server.UpdateOptionQuestion)

	// Websocket
	authRoutes.GET("/ws/main-user/:app/:room_type", func(ctx *gin.Context) {
		roomType := ctx.Param("room_type")
		log.Println("room type:", roomType)
		app := ctx.Param("app")
		//Context = ctx
		result := val.ContainRoomType(roomType)
		if !result {
			err := fmt.Errorf("this path does not exist")
			ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
		}
		user, err := HandleGetUser(ctx, server)
		if err != nil {
			err := fmt.Errorf("this path does not exist")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		}
		ServeWs(ctx.Writer, ctx.Request, roomType, user.ID, user.UserID, app, server, ctx, user)
	})

	server.router = router
}
