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

func (server *Server) RemoveGeneralReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveGeneralReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemoveGeneralReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.RemoveChargeReview(ctx, chargeID)
	if err != nil {
		log.Printf("There an error at RemoveGeneralReview at RemoveChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("your review could not be updated. Please keep in mind that you will be unable to make changes to your review 14 days after your checkout date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ReviewRes{
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   "none",
		PreviousState:  "none",
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveDetailReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemoveDetailReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		Environment: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		Accuracy: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		CheckIn: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		Communication: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		Location: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
		CurrentState: pgtype.Text{
			String: utils.DetailReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.GeneralReview,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: "started",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at RemoveDetailReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateGeneralReviewParams{
		General:        int(review.General),
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemovePrivateNoteReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemovePrivateNoteReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
		CurrentState: pgtype.Text{
			String: utils.PrivateNoteReview,
			Valid:  true,
		},
		PreviousState: pgtype.Text{
			String: utils.DetailReview,
			Valid:  true,
		},
		PrivateNote: pgtype.Text{
			String: "none",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at RemovePrivateNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateDetailReviewParams{
		Environment:    int(review.Environment),
		Accuracy:       int(review.Accuracy),
		CheckIn:        int(review.CheckIn),
		Communication:  int(review.Communication),
		Location:       int(review.Location),
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemovePublicNoteReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemovePublicNoteReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var review db.ChargeReview
	switch req.MainOptionType {
	case "events":
		review, err = server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
			PublicNote: pgtype.Text{
				String: "none",
				Valid:  true,
			},
			Status: pgtype.Text{
				String: "filled",
				Valid:  true,
			},
			CurrentState: pgtype.Text{
				String: utils.PublicNoteReview,
				Valid:  true,
			},
			PreviousState: pgtype.Text{
				String: utils.PrivateNoteReview,
				Valid:  true,
			},
			ChargeID: chargeID,
		})

		if err != nil {
			log.Printf("There an error at RemovePublicNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
			err = fmt.Errorf("could not update your review")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "options":
		review, err = server.store.UpdateChargeReview(ctx, db.UpdateChargeReviewParams{
			PublicNote: pgtype.Text{
				String: "none",
				Valid:  true,
			},
			Status: pgtype.Text{
				String: "filled",
				Valid:  true,
			},
			CurrentState: pgtype.Text{
				String: utils.PublicNoteReview,
				Valid:  true,
			},
			PreviousState: pgtype.Text{
				String: utils.PrivateNoteReview,
				Valid:  true,
			},
			ChargeID: chargeID,
		})

		if err != nil {
			log.Printf("There an error at RemovePublicNoteReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
			err = fmt.Errorf("could not update your review")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

	}

	res := CreatePrivateNoteReviewParams{
		PrivateNote:    review.PrivateNote,
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveStayCleanReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemoveStayCleanReview")

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
			String: utils.PublicNoteReview,
			Valid:  true,
		},
		StayClean: pgtype.Text{
			String: "none",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at RemoveStayCleanReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreatePublicNoteReviewParams{
		PublicNote:     review.PublicNote,
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveComfortReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemoveComfortReview")

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
			String: utils.StayCleanReview,
			Valid:  true,
		},
		StayComfort: pgtype.Text{
			String: "none",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at RemoveComfortReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateStayCleanReviewParams{
		StayClean:      review.StayClean,
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveHostReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, _, _, user, err := HandleGetCharge(ctx, server, req.ChargeID, req.MainOptionType, "RemoveHostReview")

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
		HostReview: pgtype.Text{
			String: "none",
			Valid:  true,
		},
		ChargeID: chargeID,
	})

	if err != nil {
		log.Printf("There an error at RemoveHostReview at UpdateChargeReview: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateComfortReviewParams{
		StayComfort:    review.StayComfort,
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveAllAmenityReview(ctx *gin.Context) {
	var req RemoveReviewParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveReviewParams in ShouldBindJSON: %v, ChargeID: %v \n", err.Error(), req.ChargeID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, chargeID, err := HandleOptionGetCharge(ctx, server, req.ChargeID, "options", "RemoveAllAmenityReview")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	review, err := server.store.UpdateChargeReviewAmenitiesTwo(ctx, db.UpdateChargeReviewAmenitiesTwoParams{
		Amenities:     []string{""},
		CurrentState:  utils.AmenitiesReview,
		PreviousState: utils.StayHostReview,
		ChargeID:      chargeID,
	})
	if err != nil {
		log.Printf("There an error at RemoveAllAmenityReview at UpdateChargeReviewAmenitiesTwo: %v, chargeID: %v, userID: %v \n", err.Error(), chargeID, user.ID)
		err = fmt.Errorf("could not update your review")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CreateHostReviewParams{
		HostReview:     review.HostReview,
		ChargeID:       tools.UuidToString(chargeID),
		MainOptionType: req.MainOptionType,
		CurrentState:   review.CurrentState,
		PreviousState:  review.PreviousState,
	}
	ctx.JSON(http.StatusOK, res)
}
