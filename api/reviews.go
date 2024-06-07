package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) GetStateReview(ctx *gin.Context) {
	var req GetStateReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetStateReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "GetStateReview")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	review, err := server.store.GetChargeReview(ctx, chargeID)
	if err != nil {
		log.Printf("There an error at GetStateReview at GetChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		if err == db.ErrorRecordNotFound {
			res := ReviewRes{
				ChargeID:       tools.UuidToString(chargeID),
				MainOptionType: req.MainOptionType,
				CurrentState:   utils.GeneralReview,
				PreviousState:  utils.GeneralReview,
			}
			ctx.JSON(http.StatusOK, res)
		} else {
			err = fmt.Errorf("your review could not be found")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
		}
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateGeneralReview(ctx *gin.Context) {
	var req CreateGeneralReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateGeneralReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, chargeType, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreateGeneralReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// First we remove it if it is exist
	err = server.store.RemoveChargeReview(ctx, chargeID)
	if err != nil {
		log.Printf("There an error at CreateGeneralReview at RemoveChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		if err == db.ErrorRecordNotFound {
			err = nil
		} else {
			err = fmt.Errorf("could not setup your review states")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	review, err := server.store.CreateChargeReview(ctx, db.CreateChargeReviewParams{
		ChargeID:      chargeID,
		Type:          chargeType,
		General:       int32(req.General),
		Amenities:     []string{"none"},
		CurrentState:  utils.DetailReview,
		PreviousState: utils.GeneralReview,
	})

	if err != nil {
		log.Printf("There an error at CreateGeneralReview at CreateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not setup your review, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateDetailReview(ctx *gin.Context) {
	var req CreateDetailReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateDetailReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreateDetailReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		Environment: pgtype.Int4{
			Int32: int32(req.Environment),
			Valid: true,
		},
		Accuracy: pgtype.Int4{
			Int32: int32(req.Accuracy),
			Valid: true,
		},
		CheckIn: pgtype.Int4{
			Int32: int32(req.CheckIn),
			Valid: true,
		},
		Communication: pgtype.Int4{
			Int32: int32(req.Communication),
			Valid: true,
		},
		Location: pgtype.Int4{
			Int32: int32(req.Location),
			Valid: true,
		},
		CurrentState: pgtype.Text{
			String: utils.PrivateNoteReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.DetailReview,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: "filled",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreateDetailReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreatePrivateNoteReview(ctx *gin.Context) {
	var req CreatePrivateNoteReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreatePrivateNoteReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreatePrivateNoteReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.PublicNoteReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.PrivateNoteReview,
			Valid:  true,
		},
		PrivateNote: pgtype.Text{
			String: req.PrivateNote,
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreatePrivateNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreatePublicNoteReview(ctx *gin.Context) {
	var req CreatePublicNoteReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreatePublicNoteReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreatePublicNoteReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var review db.ChargeReview
	switch req.MainOptionType {
	case "events":
		review, err = server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
			PublicNote: pgtype.Text{
				String: req.PublicNote,
				Valid:  true,
			},
			Status: pgtype.Text{
				String: "completed",
				Valid:  true,
			},
			IsPublished: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			CurrentState: pgtype.Text{
				String: utils.EventReview,
				Valid:  true,
			},
			PreviousState: pgtype.Text{
				String: utils.PublicNoteReview,
				Valid:  true,
			},
			ChargeID: chargeID,
		})

		if err != nil {
			log.Printf("There an error at CreatePublicNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
			err = fmt.Errorf("could not update your review")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "options":
		review, err = server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
			PublicNote: pgtype.Text{
				String: req.PublicNote,
				Valid:  true,
			},
			Status: pgtype.Text{
				String: "pre_complete",
				Valid:  true,
			},
			CurrentState: pgtype.Text{
				String: utils.ExtraReviewPlaceholder,
				Valid:  true,
			},
			PreviousState: pgtype.Text{
				String: utils.PublicNoteReview,
				Valid:  true,
			},
			ChargeID: chargeID,
		})

		if err != nil {
			log.Printf("There an error at CreatePublicNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
			err = fmt.Errorf("could not update your review")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

	}

	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) PlaceholderReview(ctx *gin.Context) {
	var req PlaceholderReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreatePublicNoteReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreatePublicNoteReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.StayCleanReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.ExtraReviewPlaceholder,
			Valid:  true,
		},
		ChargeID: chargeID,
	})
	if err != nil {
		log.Printf("There an error at CreatePublicNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateStayCleanReview(ctx *gin.Context) {
	var req CreateStayCleanReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateStayCleanReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreateStayCleanReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.StayComfortReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.PublicNoteReview,
			Valid:  true,
		},
		StayClean: pgtype.Text{
			String: req.StayClean,
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreateStayCleanReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateComfortReview(ctx *gin.Context) {
	var req CreateComfortReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateComfortReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreateComfortReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.StayHostReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.StayComfortReview,
			Valid:  true,
		},
		StayComfort: pgtype.Text{
			String: req.StayComfort,
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreateComfortReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateHostReview(ctx *gin.Context) {
	var req CreateHostReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateHostReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "CreateHostReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.AmenitiesReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.StayHostReview,
			Valid:  true,
		},
		HostReview: pgtype.Text{
			String: req.HostReview,
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreateHostReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateAmenityReview(ctx *gin.Context) {
	var req CreateAmenityReviewItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateAmenityReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, optionReview, chargeID, err := HandleOptionGetCharge(ctx, server, req.ChargeID, "options", "CreateAmenityReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We check if the tag already exist
	result := HandleAddAmenityReviewItem(optionReview.Amenities, req)

	review, err := server.store.UpdateChargeReviewAmenities(ctx, db.UpdateChargeReviewAmenitiesParams{
		Amenities: result,
		ChargeID:  chargeID,
	})
	if err != nil {
		log.Printf("There an error at CreateAmenityReview at UpdateChargeReviewAmenities: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := HandleAmenityReviewRes(review.Amenities, chargeID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListAmenityReview(ctx *gin.Context) {
	var req ListAmenityReviewItemParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListAmenityReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, chargeID, err := HandleOptionGetCharge(ctx, server, req.ChargeID, "options", "ListAmenityReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	emptyData := []string{"none"}
	var isEmpty bool

	//Get Amenities in the snap shot
	amenities, err := server.store.GetOptionReferenceInfoAmenities(ctx, chargeID)
	if err != nil || len(amenities) == 0 {
		if err != nil {
			log.Printf("There an error at ListAmenityReview at GetOptionReferenceInfoAmenities: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		}
		amenities = emptyData
		isEmpty = true
		err = nil
	}
	amReviews, err := server.store.GetChargeReviewAmenities(ctx, chargeID)
	if err != nil {
		log.Printf("There an error at ListAmenityReview at GetOptionReferenceInfoAmenities: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = nil
	}
	resReviews := HandleAmenityReviewRes(amReviews, chargeID)
	res := ListAmenityReviewItemRes{
		Amenities: amenities,
		IsEmpty:   isEmpty,
		Selected:  resReviews,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveAmenityReview(ctx *gin.Context) {
	var req RemoveAmenityReviewItem
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveAmenityReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, optionReview, chargeID, err := HandleOptionGetCharge(ctx, server, req.ChargeID, "options", "RemoveAmenityReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We check if the tag already exist
	result := HandleRemoveAmenityReviewItem(optionReview.Amenities, req)

	review, err := server.store.UpdateChargeReviewAmenities(ctx, db.UpdateChargeReviewAmenitiesParams{
		Amenities: result,
		ChargeID:  chargeID,
	})
	if err != nil {
		log.Printf("There an error at RemoveAmenityReview at UpdateChargeReviewAmenities: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := HandleAmenityReviewRes(review.Amenities, chargeID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CompleteOptionReview(ctx *gin.Context) {
	var req CompleteOptionReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CompleteOptionReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, "options", "CompleteOptionReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.OptionReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.AmenitiesReview,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: "completed",
			Valid:  true,
		},
		IsPublished: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at CreateHostReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: "options",
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}
