package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) CreateWishlist(ctx *gin.Context) {
	var req CreateWishlistParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateWishlist in ShouldBindJSON: %v, Selection: %v \n", err.Error(), req.OptionUserID)
		err = fmt.Errorf("request was wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  CreateWishlist in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	//  check if that optionID exist
	optionUserID, err = server.store.GetOptionForWishlist(ctx, db.GetOptionForWishlistParams{
		OptionUserID: optionUserID,
		IsActive:     true,
		IsComplete:   true,
		IsActive_2:   true,
	})
	if err != nil {
		log.Printf("Error at  CreateWishlist in GetOptionForWishlist err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item is inactive or has been unlisted")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// We check if this wishlist group already has this optionUserID
	_, err = server.store.GetWishlistItemByOptionID(ctx, db.GetWishlistItemByOptionIDParams{
		OptionUserID: optionUserID,
		UserID:       user.ID,
	})
	if err == nil {
		// This means this option is already in this wishlist group
		err = fmt.Errorf("this item already exist in this wishlist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetWishlistByName(ctx, db.GetWishlistByNameParams{
		Name:   strings.ToLower(req.Name),
		UserID: user.ID,
	})
	if err == nil {
		// We send an error because we expect that there should be an error saying it doesn't exist
		err = fmt.Errorf("this name has already been used by you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	} else {
		log.Printf("Error at  CreateWishlist in GetWishlistByName err: %v, user: %v\n", err, user.ID)
	}
	wishlist, err := server.store.CreateWishlist(ctx, db.CreateWishlistParams{
		Name:   req.Name,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at  CreateWishlist in CreateWishlistItem err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not create your wishlist %v", req.Name)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Next we create the wishlist item
	wishlistItem, err := server.store.CreateWishlistItem(ctx, db.CreateWishlistItemParams{
		WishlistID:   wishlist.ID,
		OptionUserID: optionUserID,
	})
	if err != nil {
		log.Printf("Error at  CreateWishlist in CreateWishlistItem err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not add item to %v", req.Name)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	mainImage := "none"
	optionPhotos, err := server.store.GetOptionInfoPhotoByOptionUserID(ctx, wishlistItem.OptionUserID)
	if err != nil {
		log.Printf("Error at  CreateWishlist in GetOptionInfoPhotoByOptionUserID err: %v, user: %v\n", err, user.ID)
	} else {
		mainImage = optionPhotos.MainImage
	}
	res := WishlistItem{
		Name:           wishlist.Name,
		WishlistID:     tools.UuidToString(wishlist.ID),
		WishlistItemID: tools.UuidToString(wishlistItem.ID),
		OptionUserID:   tools.UuidToString(wishlistItem.OptionUserID),
		MainImage:     mainImage,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateWishlistItem(ctx *gin.Context) {
	var req CreateWishlistItemParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateWishlistItem in ShouldBindJSON: %v, Selection: %v \n", err.Error(), req.OptionUserID)
		err = fmt.Errorf("request was wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  CreateWishlistItem optionUserID in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	//  check if that optionID exist
	optionUserID, err = server.store.GetOptionForWishlist(ctx, db.GetOptionForWishlistParams{
		OptionUserID: optionUserID,
		IsActive:     true,
		IsComplete:   true,
		IsActive_2:   true,
	})
	if err != nil {
		log.Printf("Error at CreateWishlistItem in GetOptionForWishlist err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item is inactive or has been unlisted")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	wishlistID, err := tools.StringToUuid(req.WishlistID)
	if err != nil {
		log.Printf("Error at  CreateWishlistItem wishlistID in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	wishlist, err := server.store.GetWishlist(ctx, db.GetWishlistParams{
		ID:     wishlistID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at  CreateWishlistItem in GetWishlist err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("this wishlist group cannot be found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// We check if this wishlist group already has this optionUserID
	_, err = server.store.GetWishlistItemByOptionID(ctx, db.GetWishlistItemByOptionIDParams{
		OptionUserID: optionUserID,
		UserID:       user.ID,
	})
	if err == nil {
		// This means this option is already in this wishlist group
		err = fmt.Errorf("this item already exist in this wishlist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Next we create the wishlist item
	wishlistItem, err := server.store.CreateWishlistItem(ctx, db.CreateWishlistItemParams{
		WishlistID:   wishlist.ID,
		OptionUserID: optionUserID,
	})
	if err != nil {
		log.Printf("Error at  CreateWishlistItem in CreateWishlistItem err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not add item to this wishlist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	mainImage := "none"
	optionPhotos, err := server.store.GetOptionInfoPhotoByOptionUserID(ctx, wishlistItem.OptionUserID)
	if err != nil {
		log.Printf("Error at  CreateWishlist in GetOptionInfoPhotoByOptionUserID err: %v, user: %v\n", err, user.ID)
	} else {
		mainImage = optionPhotos.MainImage
	}
	res := WishlistItem{
		Name:           wishlist.Name,
		WishlistID:     tools.UuidToString(wishlist.ID),
		WishlistItemID: tools.UuidToString(wishlistItem.ID),
		OptionUserID:   tools.UuidToString(wishlistItem.OptionUserID),
		MainImage:     mainImage,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveWishlist(ctx *gin.Context) {
	var req RemoveWishlistParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveWishlist in ShouldBindJSON: %v, Selection: %v \n", err.Error(), req.WishlistID)
		err = fmt.Errorf("request was wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	wishlistID, err := tools.StringToUuid(req.WishlistID)
	if err != nil {
		log.Printf("Error at  RemoveWishlist wishlistID in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	wishlistItems, err := server.store.ListWishlistItem(ctx, db.ListWishlistItemParams{
		WishlistID: wishlistID,
		UserID:     user.ID,
	})
	if err != nil {
		log.Printf("Error at RemoveWishlist wishlistID in ListWishlistItem err: %v, user: %v\n", err, user.ID)
		err = errors.New("an error occurred while removing your wishlist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// We want to remove all wishlist items
	for _, item := range wishlistItems {
		err := server.store.RemoveWishlistItem(ctx, item)
		if err != nil {
			log.Printf("Error at RemoveWishlist wishlistID in L.RemoveWishlistItem err: %v, user: %v\n", err, user.ID)
			err = errors.New("an error occurred while removing your wishlist")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}
	// We have remove all the wishlist items
	// We want to remove the main wishlist
	err = server.store.RemoveWishlist(ctx, db.RemoveWishlistParams{
		ID:     wishlistID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at RemoveWishlist wishlistID in L.RemoveWishlistItem err: %v, user: %v\n", err, user.ID)
		err = errors.New("an error occurred while removing your wishlist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveWishlistItem(ctx *gin.Context) {
	var req RemoveWishlistItemParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveWishlistItem in ShouldBindJSON: %v, Selection: %v \n", err.Error(), req.WishlistItemID)
		err = fmt.Errorf("request was wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wishlistItemID, err := tools.StringToUuid(req.WishlistItemID)
	if err != nil {
		log.Printf("Error at  RemoveWishlistItem wishlistItemID in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	wishlistItem, err := server.store.GetWishlistItem(ctx, db.GetWishlistItemParams{
		ID:     wishlistItemID,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error at  RemoveWishlistItem GetWishlistItem in StringToUuid err: %v, user: %v\n", err, user.ID)
		err = errors.New("this item does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// We want to remove the given wishlist item
	err = server.store.RemoveWishlistItem(ctx, wishlistItemID)
	if err != nil {
		log.Printf("Error at RemoveWishlistItem wishlistID in RemoveWishlistItem err: %v, user: %v\n", err, user.ID)
		err = errors.New("an error occurred while removing your wishlist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// We check to see how the album is if there are anymore wishlist
	wishlistCount, err := server.store.GetWishlistItemCountAll(ctx, db.GetWishlistItemCountAllParams{
		WishlistID: wishlistItem.WishlistID,
		UserID:     user.ID,
	})
	if err != nil {
		log.Printf("Error at RemoveWishlistItem wishlistID in store.GetWishlistItemCountAll err: %v, user: %v\n", err, user.ID)
	} else {
		if wishlistCount < 1 {
			// If less than one means it is empty so we want to delete it
			err = server.store.RemoveWishlist(ctx, db.RemoveWishlistParams{
				UserID: user.ID,
				ID:     wishlistItem.WishlistID,
			})
			if err != nil {
				log.Printf("Error at RemoveWishlistItem wishlistID in RemoveWishlist err: %v, user: %v\n", err, user.ID)
			}
		}
	}
	res := RemoveWishlistItemRes{
		Success:      true,
		OptionUserID: tools.UuidToString(wishlistItem.OptionUserID),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListWishlist(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := HandleWishlist(ctx, server, user)

	log.Printf("user logged in successfully (%v) \n", user.Email)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListWishlistExperience(ctx *gin.Context) {
	var req WishlistOffsetParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListExperience in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionOffset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var resOption ListExperienceWishlistOptionRes
	var resEvent ListExperienceWishlistEventRes
	var hasData bool
	switch req.MainOptionType {
	case "options":
		resOption, err, hasData = HandleWishlistOptionExperience(ctx, server, user, req)

	case "events":
		resEvent, err, hasData = HandleWishlistEventExperience(ctx, server, user, req)
	}
	if hasData && err == nil {
		switch req.MainOptionType {
		case "options":
			ctx.JSON(http.StatusOK, resOption)
			return
		case "events":
			ctx.JSON(http.StatusOK, resEvent)
			return
		}
	} else if !hasData && err == nil {
		ctx.JSON(http.StatusNoContent, "none")
		return
	} else {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
