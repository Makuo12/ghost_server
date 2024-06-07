package api

import (
	"log"
	"strings"

	"github.com/makuo12/ghost_server/algo"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetOptionCategories(category db.OptionsInfosCategory) (result []string) {
	data := make(map[string]int64)
	list := [][]string{
		category.TypeOfShortlet,
		category.Amenities,
		category.SpaceType,
		category.Des,
		category.Name,
		category.Highlight,
		category.SpaceArea}
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

	maxInt := int64(0)
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

func CreateOptionAlgo(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {

	var des string
	var name string
	var amenities []string
	var spaceAreas []string
	var highlights []string
	var shortletType string
	var spaceType string

	// Des
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error HandleAllOptionAlgo at GetOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		des = ""
		name = ""
		highlights = []string{""}
	} else {
		des = tools.HandleString(optionDetail.Des)
		name = tools.HandleString(optionDetail.HostNameOption)
		highlights = tools.HandleDBList(optionDetail.OptionHighlight)
	}

	desRatio := algo.HandleAlgoData(algo.HandleOptionAlgoDes(des))
	nameRatio := algo.HandleAlgoData(algo.HandleOptionAlgoName(name))
	highlightRatio := algo.HandleAlgoData(algo.HandleOptionAlgoHigh(highlights))

	// Amenities
	amenitiesData, err := server.store.ListAmenitiesTag(ctx, db.ListAmenitiesTagParams{
		OptionID: option.ID,
		HasAm:    true,
	})
	if err != nil {
		log.Printf("There an error HandleAllOptionAlgo at ListAmenitiesTag: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		amenities = []string{""}
	} else {
		amenities = amenitiesData
	}
	amenitiesRatio := algo.HandleAlgoData(algo.HandleOptionAlgoAmenities(amenities))

	// SpaceAreas
	spaceAreaData, err := server.store.ListSpaceAreaType(ctx, option.ID)
	if err != nil {
		log.Printf("There an error HandleAllOptionAlgo at ListSpaceAreaType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		spaceAreas = []string{""}
	} else {
		spaceAreas = spaceAreaData
	}
	spaceAreasRatio := algo.HandleAlgoData(algo.HandleOptionAlgoSpaceAreas(spaceAreas))

	// Shortlet
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("There an error HandleAllOptionAlgo at GetShortlet: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		shortletType = ""
		spaceType = ""
	} else {
		shortletType = shortlet.TypeOfShortlet
		spaceType = shortlet.SpaceType
	}
	shortletTypeRatio := algo.HandleAlgoData(algo.HandleOptionAlgoType(shortletType))
	spaceTypeRatio := algo.HandleAlgoData(algo.HandleOptionAlgoSpaceType(spaceType))

	category, err := server.store.CreateOptionInfoCategory(ctx, db.CreateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: shortletTypeRatio,
		SpaceArea:      spaceAreasRatio,
		SpaceType:      spaceTypeRatio,
		Highlight:      highlightRatio,
		Des:            desRatio,
		Name:           nameRatio,
		Amenities:      amenitiesRatio,
	})

	if err != nil {
		log.Printf("There an error HandleAllOptionAlgo at CreateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error HandleAllOptionAlgo at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategoryDes(des string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryDes at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	desRatio := algo.HandleAlgoData(algo.HandleOptionAlgoDes(des))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      categoryData.SpaceType,
		Highlight:      categoryData.Highlight,
		Des:            desRatio,
		Name:           categoryData.Name,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryDes at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategoryDes at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategoryName(name string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryName at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	nameRatio := algo.HandleAlgoData(algo.HandleOptionAlgoName(name))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      categoryData.SpaceType,
		Highlight:      categoryData.Highlight,
		Des:            categoryData.Des,
		Name:           nameRatio,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryName at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategoryName at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategoryAmenities(arr []string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryAmenities at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	amenitiesRatio := algo.HandleAlgoData(algo.HandleOptionAlgoAmenities(arr))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      categoryData.SpaceType,
		Highlight:      categoryData.Highlight,
		Des:            categoryData.Des,
		Name:           categoryData.Name,
		Amenities:      amenitiesRatio,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryAmenities at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategoryAmenities at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategorySpaceAreas(arr []string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategorySpaceAreas at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	spaceAreaRatio := algo.HandleAlgoData(algo.HandleOptionAlgoSpaceAreas(arr))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      spaceAreaRatio,
		SpaceType:      categoryData.SpaceType,
		Highlight:      categoryData.Highlight,
		Des:            categoryData.Des,
		Name:           categoryData.Name,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategorySpaceAreas at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategorySpaceArea at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategoryHigh(arr []string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryHigh at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	highlightRatio := algo.HandleAlgoData(algo.HandleOptionAlgoHigh(arr))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      categoryData.SpaceType,
		Highlight:      highlightRatio,
		Des:            categoryData.Des,
		Name:           categoryData.Name,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryHigh at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategoryHigh at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategoryType(value string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryType at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	typeRatio := algo.HandleAlgoData(algo.HandleOptionAlgoType(value))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: typeRatio,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      categoryData.SpaceType,
		Highlight:      categoryData.Highlight,
		Des:            categoryData.Des,
		Name:           categoryData.Name,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategoryType at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategoryType at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}

func UpdateOptionCategorySpaceType(value string, ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User) {
	categoryData, err := server.store.GetOptionInfoCategory(ctx, option.ID)
	if err != nil {
		log.Printf("There an error UpdateOptionCategorySpaceType at GetOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	typeRatio := algo.HandleAlgoData(algo.HandleOptionAlgoSpaceType(value))
	category, err := server.store.UpdateOptionInfoCategory(ctx, db.UpdateOptionInfoCategoryParams{
		OptionID:       option.ID,
		TypeOfShortlet: categoryData.TypeOfShortlet,
		SpaceArea:      categoryData.SpaceArea,
		SpaceType:      typeRatio,
		Highlight:      categoryData.Highlight,
		Des:            categoryData.Des,
		Name:           categoryData.Name,
		Amenities:      categoryData.Amenities,
	})
	if err != nil {
		log.Printf("There an error UpdateOptionCategorySpaceType at UpdateOptionInfoCategory: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		return
	}
	result := GetOptionCategories(category)
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
		log.Printf("There an error UpdateOptionCategorySpaceType at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

}
