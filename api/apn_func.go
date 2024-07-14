package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/makuo12/ghost_server/constants"
	"github.com/makuo12/ghost_server/tools"

	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
)

type GetAppVersionRes struct {
	Channel string `json:"channel"`
	Msg     string `json:"msg"`
}

func (server *Server) GetAppVersion(ctx *gin.Context) {
	appType := ctx.Param("type")
	version := ctx.Param("version")
	var channel string
	var msg string
	switch appType {
	case "ios_app":
		if version == server.config.IOSAppVersion {
			channel = constants.VERSION_STABLE
			msg = "none"
		} else {
			channel = constants.VERSION_UNSTABLE
			msg = server.config.IOSAppVersionMsg
		}

	case "android_app":
		err := fmt.Errorf("app not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	res := GetAppVersionRes{
		Channel: channel,
		Msg:     msg,
	}
	ctx.JSON(http.StatusOK, res)
}

func HandleUIDApn(ctx context.Context, server *Server, uID uuid.UUID, title string, msg string) {
	apns, err := server.store.ListUidAPNDetail(ctx, uID)
	if err != nil {
		log.Println("err at HandleUIDApn apn: ", err.Error(), "for uID: ", uID)
		return
	}
	for _, apn := range apns {
		SendApn(ctx, server, apn.Token, msg, title)
	}
}

func HandleUserIdApn(ctx context.Context, server *Server, userID uuid.UUID, title string, msg string) {
	apns, err := server.store.ListUserIdAPNDetail(ctx, userID)
	if err != nil {
		log.Println("err at HandleUserIdApn apn: ", err.Error(), "for userID: ", userID)
		return
	}
	for _, apn := range apns {
		SendApn(ctx, server, apn.Token, msg, title)
	}
}

func HandleUserIdMessageApn(ctx context.Context, server *Server, userID uuid.UUID, msg string, name string) {
	apns, err := server.store.ListUserIdAPNDetail(ctx, userID)
	if err != nil {
		log.Println("err at HandleUserIdApn apn: ", err.Error(), "for userID: ", userID)
		return
	}

	for _, apn := range apns {
		title := fmt.Sprintf("Message from %v", tools.CapitalizeFirstCharacter(name))
		SendApn(ctx, server, apn.Token, msg, title)
	}
}

func SendApn(ctx context.Context, server *Server, deviceToken string, msg string, title string) {
	// Obtain a messaging.Client from the App.

	// This registration token comes from the client FCM SDKs.
	registrationToken := deviceToken
	badgeNumber := 1
	// See documentation on defining a message payload.
	notification := &messaging.Message{
    Notification: &messaging.Notification{
        Title: title,
        Body:  msg,
    },
    APNS: &messaging.APNSConfig{
        Headers: map[string]string{
            "apns-priority": "10",
        },
        Payload: &messaging.APNSPayload{
            Aps: &messaging.Aps{
                Sound: "default",
                Badge: &badgeNumber, // Add your badge number here
            },
        },
    },
    Token: registrationToken,
}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := server.ApnFire.Send(ctx, notification)
	if err != nil {
		log.Println("err at send apn: ", err.Error())
		return
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}
