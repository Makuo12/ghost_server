package api

import (
	"context"
	"flex_server/algo"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func IsExperienceEmpty() bool {
	result, err := RedisClient.SMembers(RedisContext, constants.ALL_EXPERIENCE_CATEGORIES).Result()
	if err != nil || len(result) == 0 {
		if err != nil {
			log.Printf("Error at HandleRedisOptionExperience in RedisClient.SMembers(RedisContext err: %v\n", err)
		}
		return false
	}
	return true
}

func HandleRedisOptionExperience(ctx *gin.Context, server *Server, req ExperienceOffsetParams) (res ListExperienceOptionRes, err error, hasData bool) {
	hasData = true
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	catKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_OPTION, req.Type)
	result, err := RedisClient.SMembers(RedisContext, catKey).Result()
	if err != nil {
		log.Printf("Error at HandleRedisOptionExperience in RedisClient.SMembers(RedisContext err: %v, user: %v, catKey: %v\n", err, ctx.ClientIP(), catKey)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if len(result) <= req.OptionOffset {
		err = nil
		hasData = false
		return
	}
	var resData []ExperienceOptionData
	for _, r := range result {
		data, err := RedisClient.HGetAll(RedisContext, r).Result()
		if err != nil {
			log.Printf("Error at HandleRedisOptionExperience in RedisClient.HGetAll( err: %v, user: %v, catKey: %v id: %v\n", err, ctx.ClientIP(), catKey, r)
			continue
		}
		optionUserID, err := tools.StringToUuid(getDataValue(data, constants.OPTION_USER_ID))
		if err != nil {
			log.Printf("Error at HandleRedisOptionExperience in StringToUuid(getDataValue(data, constants.OPTION_USER_ID) iD: %v err: %v, user: %v, catKey: %v id: %v\n", getDataValue(data, constants.OPTION_USER_ID), err, ctx.ClientIP(), catKey, r)
			continue
		}
		newData := ConvertSliceToExperienceOptionData(data, dollarToNaira, dollarToCAD, req.Currency, optionUserID)
		resData = append(resData, newData)
	}
	if len(resData) == 0 {
		hasData = false
		err = nil
		return
	}
	resData = CustomSort(resData, req.State, req.Country)
	log.Println("resData redis experience", resData)
	resIndexData := GetExperienceOptionOffset(resData, req.OptionOffset, 10)
	log.Println("resIndexData redis experience", resIndexData)
	if err == nil && hasData {
		res = ListExperienceOptionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(resIndexData),
			OnLastIndex:  false,
			Category:     req.Type,
		}
	}
	return
}

func HandleRedisEventExperience(ctx *gin.Context, server *Server, req ExperienceOffsetParams) (res ListExperienceEventRes, err error, hasData bool) {
	hasData = true
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	catKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_EVENT, req.Type)
	log.Println("catKey experience", catKey)
	result, err := RedisClient.SMembers(RedisContext, catKey).Result()
	if err != nil {
		log.Printf("Error at HandleRedisEventExperience in RedisClient.SMembers(RedisContext err: %v, user: %v, catKey: %v\n", err, ctx.ClientIP(), catKey)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	log.Println("result experience", result)
	if len(result) <= req.OptionOffset {
		err = nil
		hasData = false
		return
	}
	var resData []ExperienceEventData
	for _, r := range result {
		data, err := RedisClient.HGetAll(RedisContext, r).Result()
		if err != nil {
			log.Printf("Error at HandleRedisEventExperience in RedisClient.HGetAll( err: %v, user: %v, catKey: %v id: %v\n", err, ctx.ClientIP(), catKey, r)
			continue
		}
		log.Println("data exr", data)
		optionUserID, err := tools.StringToUuid(getDataValue(data, constants.OPTION_USER_ID))
		if err != nil {
			log.Printf("Error at HandleRedisEventExperience in StringToUuid(getDataValue(data, constants.OPTION_USER_ID) iD: %v err: %v, user: %v, catKey: %v id: %v\n", getDataValue(data, constants.OPTION_USER_ID), err, ctx.ClientIP(), catKey, r)
			continue
		}
		newData := ConvertSliceToExperienceEventData(data, dollarToNaira, dollarToCAD, req.Currency, optionUserID)
		resData = append(resData, newData)
	}
	if len(resData) == 0 {
		hasData = false
		err = nil
		return
	}
	resData = CustomEventExperienceDataSort(resData, req.State, req.Country)
	resIndexData := GetExperienceEventOffset(resData, req.OptionOffset, 10)
	log.Println("resIndexData redis experience", resIndexData)
	if err == nil && hasData {
		res = ListExperienceEventRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(resIndexData),
			OnLastIndex:  false,
			Category:     req.Type,
		}
	}
	return
}

func HandleEventExperienceToRedis(ctx context.Context, server *Server) func() {
	return func() {
		log.Println("Here at redis setup event experience")
		var catKeys []string
		for _, cat := range algo.EventCategory {
			catKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_EVENT, cat)
			eventInfos, err := server.store.ListEventExperience(ctx, db.ListEventExperienceParams{
				IsComplete:     true,
				IsActive:       true,
				IsActive_2:     true,
				MainOptionType: "events",
			})
			if err != nil || len(eventInfos) == 0 {
				if err != nil {
					log.Printf("Error at HandleEventExperienceToRedis in ListEventExperience(ctx, err: %v\n", err)
				}
				continue
			}
			var eventKeys []string
			for _, e := range eventInfos {
				eventKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_EVENT, e.OptionUserID)
				cats := []string{e.Category, e.CategoryTwo, e.CategoryThree}
				if e.Status == "unlist" || e.Status == "snooze" || !tools.IsInList(cats, cat) {
					if e.Status == "unlist" || e.Status == "snooze" {
						err := RedisClient.Del(RedisContext, eventKey).Err()
						if err != nil {
							log.Printf("Error at HandleOptionExperienceToRedis in unlist in RedisClient.Del, eventKey %v err: %v\n", eventKey, err)
						}
					}
					// We want to remove it from the SMEMBERS
					err = RedisClient.SRem(RedisContext, catKey, eventKey).Err()
					if err != nil {
						log.Printf("Error at HandleOptionExperienceToRedis in unlist in RedisClient.SRem, eventKey %v err: %v\n", eventKey, err)
					}
					continue
				}
				locationRedisList, _, price, ticketAvailable, startDateData, endDateData, hasFreeTicket := SetupExperienceEventData(ctx, server, e, db.ListEventExperienceByLocationRow{}, true, "HandleEventExperienceToRedis")
				data := []string{
					constants.OPTION_USER_ID,
					tools.UuidToString(e.OptionUserID),
					constants.HOST_OPTION_NAME,
					e.HostNameOption,
					constants.OPTION_IS_VERIFIED,
					tools.ConvertBoolToString(e.IsVerified),
					constants.COVER_IMAGE,
					e.CoverImage,
					constants.HOST_AS_INDIVIDUAL,
					tools.ConvertBoolToString(e.HostAsIndividual),
					constants.PHOTOS,
					strings.Join(e.Photo, "&"),
					constants.TICKET_AVAILABLE,
					tools.ConvertBoolToString(ticketAvailable),
					constants.SUB_EVENT_TYPE,
					e.SubCategoryType,
					constants.TICKET_LOWEST_PRICE,
					tools.ConvertFloatToString(price),
					constants.EVENT_START_DATE,
					startDateData,
					constants.EVENT_END_DATE,
					endDateData,
					constants.HAS_FREE_TICKET,
					tools.ConvertBoolToString(hasFreeTicket),
					constants.HOST_NAME,
					e.FirstName,
					constants.PROFILE_PHOTO,
					e.Photo_2,
					constants.HOST_JOINED,
					tools.ConvertDateOnlyToString(e.CreatedAt),
					constants.HOST_VERIFIED,
					tools.ConvertBoolToString(e.IsVerified_2),
					constants.CATEGORY,
					e.Category,
					constants.LOCATION,
					strings.Join(locationRedisList, "&"),
					constants.CURRENCY,
					e.Currency,
				}

				err := RedisClient.HSet(RedisContext, eventKey, data).Err()
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in HSet(RedisContext, eventKey err: %v\n", err)
					continue
				}
				eventKeys = append(eventKeys, eventKey)

			}
			if len(eventKeys) != 0 {
				err = RedisClient.SAdd(RedisContext, catKey, eventKeys).Err()
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in RedisClient.SAdd(RedisContext, catKey err: %v\n", err)
					continue
				} else {
					catKeys = append(catKeys, catKey)
				}
			}
		}
		err := RedisClient.SAdd(RedisContext, constants.ALL_EXPERIENCE_CATEGORIES, catKeys).Err()
		if err != nil {
			log.Printf("Error at HandleOptionExperienceToRedis in RedisClient.SAdd(RedisContext, constants.ALL_EXPERIENCE_CATEGORIES err: %v\n", err)
		}
	}

}

func HandleOptionExperienceToRedis(ctx context.Context, server *Server) func() {
	return func() {
		log.Println("Here at redis setup options experience")
		var catKeys []string
		for _, cat := range algo.OptionCategory {
			catKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_OPTION, cat)
			optionInfos, err := server.store.ListOptionExperience(ctx, db.ListOptionExperienceParams{
				IsComplete:     true,
				IsActive:       true,
				IsActive_2:     true,
				MainOptionType: "options",
			})
			if err != nil || len(optionInfos) == 0 {
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in ListOptionExperience(ctx, err: %v\n", err)
				}
				continue
			}
			var optionKeys []string
			for _, o := range optionInfos {
				// If perform some checks
				optionKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_OPTION, o.OptionUserID)
				locationKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_OPTION_LOCATION, o.OptionUserID)
				cats := []string{o.Category, o.CategoryTwo, o.CategoryThree}
				if o.Status == "unlist" || o.Status == "snooze" || !tools.IsInList(cats, cat) {
					// We want to remove it from the SMEMBERS
					if o.Status == "unlist" || o.Status == "snooze" {
						err := RedisClient.Del(RedisContext, optionKey).Err()
						if err != nil {
							log.Printf("Error at HandleOptionExperienceToRedis in unlist in RedisClient.Del, optionKey %v err: %v\n", optionKey, err)
						}
						err = RedisClient.ZRem(RedisContext, constants.ALL_EXPERIENCE_LOCATION, locationKey).Err()
						if err != nil {
							log.Printf("Error at HandleOptionExperienceToRedis in unlist in RedisClient.ZRem(, optionKey %v err: %v\n", optionKey, err)
						}
					}
					err = RedisClient.SRem(RedisContext, catKey, optionKey).Err()
					if err != nil {
						log.Printf("Error at HandleOptionExperienceToRedis in unlist in RedisClient.SRem, optionKey %v err: %v\n", optionKey, err)
					}

					continue
				}
				data := ConvertExperienceOptionDataToSlice(o)
				err := RedisClient.HSet(RedisContext, optionKey, data).Err()
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in HSet(RedisContext, optionKey err: %v\n", err)
					continue
				}
				optionKeys = append(optionKeys, optionKey)
				lat := tools.ConvertFloatToLocationString(o.Geolocation.P.Y, 9)
				lng := tools.ConvertFloatToLocationString(o.Geolocation.P.X, 9)
				location := &redis.GeoLocation{
					Latitude:  tools.ConvertLocationStringToFloat(lat, 9),
					Longitude: tools.ConvertLocationStringToFloat(lng, 9),
					Name:      locationKey,
				}
				err = RedisClient.GeoAdd(RedisContext, constants.ALL_EXPERIENCE_LOCATION, location).Err()
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in RedisClient.GeoAdd(ctx, constants.ALL_EXPERIENCE_LOCATION err: %v\n", err)
					continue
				}
			}
			if len(optionKeys) != 0 {
				err = RedisClient.SAdd(RedisContext, catKey, optionKeys).Err()
				if err != nil {
					log.Printf("Error at HandleOptionExperienceToRedis in RedisClient.SAdd(RedisContext, catKey err: %v\n", err)
					continue
				} else {
					catKeys = append(catKeys, catKey)
				}
			}
		}
		err := RedisClient.SAdd(RedisContext, constants.ALL_EXPERIENCE_CATEGORIES, catKeys).Err()
		if err != nil {
			log.Printf("Error at HandleOptionExperienceToRedis in RedisClient.SAdd(RedisContext, constants.ALL_EXPERIENCE_CATEGORIES err: %v\n", err)
		}
	}

}
