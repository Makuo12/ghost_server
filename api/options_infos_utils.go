package api

import (
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func GetCurrentOptionInfo(server *Server, ctx *gin.Context, user db.User) (optionInfo db.OptionsInfo, err error) {
	errorCurrentOption := false
	if len(user.CurrentOptionID) != 0 && user.CurrentOptionID != "none" {
		result, err := tools.StringToUuid(user.CurrentOptionID)
		if err != nil {
			log.Printf("Error at GetCurrentOptionInfo in tools.StringToUuid: %v, userID: %v \n", err.Error(), user.ID)
			errorCurrentOption = true
		} else {
			optionInfo, err = server.store.GetOptionInfo(ctx, db.GetOptionInfoParams{
				ID:         result,
				HostID:     user.ID,
				IsComplete: true,
			})
			if err != nil {
				log.Printf("Error at GetCurrentOptionInfo in GetOptionInfo: %v, userID: %v \n", err.Error(), user.ID)
				errorCurrentOption = true
			}
		}
	} else {
		errorCurrentOption = true
	}
	// errorCurrentOption this would run if anything went wrong when fetching a current option that will feel was stored
	if errorCurrentOption {
		optionInfo, err = server.store.GetHostOptionInfo(ctx, db.GetHostOptionInfoParams{
			HostID:     user.ID,
			IsComplete: true,
		})
		if err != nil {
			log.Printf("Error at GetCurrentOptionInfo in GetHostOptionInfo: %v, userID: %v \n", err.Error(), user.ID)
		}
		return
	}
	return
}
