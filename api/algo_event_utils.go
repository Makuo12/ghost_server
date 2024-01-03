package api

import (
	"flex_server/algo"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetEventCategories(category db.EventsInfosCategory) (result []string) {
	data := make(map[string]int64)
	list := [][]string{
		category.EventType,
		category.EventSubType,
		category.Des,
		category.Name,
		category.Highlight}
	for _, catList := range list {
		for _, cat := range catList {
			item := strings.Split(cat, "&")
			value, err := tools.ConvertStringToInt64(item[1])
			if err != nil {
				value = 0
			}
			data[item[0]] = data[item[0]] + value
		}
	}

	var maxInt int64
	for _, num := range data {
		if num > maxInt {
			maxInt = num
		}
	}

	var selected []string
	for key, num := range data {
		if num == maxInt {
			selected = append(selected, key)
			// If the difference is just three it is still okay
		} else if num >= maxInt-3 {
			selected = append(selected, key)
		}
	}
	// We are doing this so we can be sure that there would always be 4
	for i := 0; i < 4; i++ {
		if i < len(selected) {
			result = append(result, selected[i])
		} else {
			result = append(result, "none")
		}
	}

	return

}

func CreateEventAlgo(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {

	var des string
	var name string
	var highlights []string
	var eventSubType string
	var eventType string

	// Des
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error HandleAllEventAlgo at GetOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		des = ""
		name = ""
		highlights = []string{""}
	} else {
		des = tools.HandleString(optionDetail.Des)
		name = tools.HandleString(optionDetail.HostNameOption)
		highlights = tools.HandleDBList(optionDetail.OptionHighlight)
	}

	desRatio := algo.HandleAlgoData(algo.HandleEventAlgoDes(des))
	nameRatio := algo.HandleAlgoData(algo.HandleEventAlgoName(name))
	highlightRatio := algo.HandleAlgoData(algo.HandleEventAlgoHigh(highlights))

	// Event

	eventInfo, err := server.store.GetEventInfo(ctx, option.ID)
	if err != nil {
		log.Printf("There an error HandleAllEventAlgo at GetEventInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		eventType = ""
		eventSubType = ""
	} else {
		eventType = eventInfo.EventType
		eventSubType = eventInfo.SubCategoryType
	}
	eventTypeRatio := algo.HandleAlgoData(algo.HandleEventAlgoType(eventType))
	eventSubTypeRatio := algo.HandleAlgoData(algo.HandleEventAlgoSubType(eventSubType))

	category, err := server.store.CreateEventInfoCategory(ctx, db.CreateEventInfoCategoryParams{
		OptionID:     option.ID,
		Highlight:    highlightRatio,
		Des:          desRatio,
		Name:         nameRatio,
		EventType:    eventTypeRatio,
		EventSubType: eventSubTypeRatio,
	})

	if err != nil {
		log.Printf("There an error HandleAllEventAlgo at CreateEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetEventCategories(category)
	_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		ID: option.ID,
		Category: pgtype.Text{
			String: result[0],
			Valid:  true,
		},
		CategoryTwo: pgtype.Text{
			String: result[1],
			Valid:  true,
		},
		CategoryThree: pgtype.Text{
			String: result[2],
			Valid:  true,
		},
		CategoryFour: pgtype.Text{
			String: result[3],
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("There an error HandleAllEventAlgo at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateEventCategoryDes(des string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetEventInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateEventCategoryDes at GetEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	desRatio := algo.HandleAlgoData(algo.HandleEventAlgoDes(des))
	category, err := server.store.UpdateEventInfoCategory(ctx, db.UpdateEventInfoCategoryParams{
		OptionID:     option.ID,
		Highlight:    categoryData.Highlight,
		Des:          desRatio,
		Name:         categoryData.Name,
		EventType:    categoryData.EventType,
		EventSubType: categoryData.EventSubType,
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryDes at UpdateEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetEventCategories(category)
	_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		ID: option.ID,
		Category: pgtype.Text{
			String: result[0],
			Valid:  true,
		},
		CategoryTwo: pgtype.Text{
			String: result[1],
			Valid:  true,
		},
		CategoryThree: pgtype.Text{
			String: result[2],
			Valid:  true,
		},
		CategoryFour: pgtype.Text{
			String: result[3],
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryDes at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateEventCategoryName(name string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetEventInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateEventCategoryName at GetEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	nameRatio := algo.HandleAlgoData(algo.HandleEventAlgoName(name))
	category, err := server.store.UpdateEventInfoCategory(ctx, db.UpdateEventInfoCategoryParams{
		OptionID:     option.ID,
		EventType:    categoryData.EventType,
		EventSubType: categoryData.EventSubType,
		Highlight:    categoryData.Highlight,
		Des:          categoryData.Des,
		Name:         nameRatio,
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryName at UpdateEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetEventCategories(category)
	_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		ID: option.ID,
		Category: pgtype.Text{
			String: result[0],
			Valid:  true,
		},
		CategoryTwo: pgtype.Text{
			String: result[1],
			Valid:  true,
		},
		CategoryThree: pgtype.Text{
			String: result[2],
			Valid:  true,
		},
		CategoryFour: pgtype.Text{
			String: result[3],
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryName at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateEventCategoryType(value string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetEventInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateEventCategoryName at GetEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	valueRatio := algo.HandleAlgoData(algo.HandleEventAlgoType(value))
	category, err := server.store.UpdateEventInfoCategory(ctx, db.UpdateEventInfoCategoryParams{
		OptionID:     option.ID,
		EventType:    valueRatio,
		EventSubType: categoryData.EventSubType,
		Highlight:    categoryData.Highlight,
		Des:          categoryData.Des,
		Name:         categoryData.Name,
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryType at UpdateEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetEventCategories(category)
	_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		ID: option.ID,
		Category: pgtype.Text{
			String: result[0],
			Valid:  true,
		},
		CategoryTwo: pgtype.Text{
			String: result[1],
			Valid:  true,
		},
		CategoryThree: pgtype.Text{
			String: result[2],
			Valid:  true,
		},
		CategoryFour: pgtype.Text{
			String: result[3],
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategoryType at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateEventCategorySubType(value string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetEventInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateEventCategoryName at GetEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	valueRatio := algo.HandleAlgoData(algo.HandleEventAlgoType(value))
	category, err := server.store.UpdateEventInfoCategory(ctx, db.UpdateEventInfoCategoryParams{
		OptionID:     option.ID,
		EventType:    categoryData.EventType,
		EventSubType: valueRatio,
		Highlight:    categoryData.Highlight,
		Des:          categoryData.Des,
		Name:         categoryData.Name,
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategorySubType at UpdateEventInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetEventCategories(category)
	_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		ID: option.ID,
		Category: pgtype.Text{
			String: result[0],
			Valid:  true,
		},
		CategoryTwo: pgtype.Text{
			String: result[1],
			Valid:  true,
		},
		CategoryThree: pgtype.Text{
			String: result[2],
			Valid:  true,
		},
		CategoryFour: pgtype.Text{
			String: result[3],
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("There an error UpdateEventCategorySubType at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}
