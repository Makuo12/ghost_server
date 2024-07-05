package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

func (server *Server) ListPhotoUserAdmin(ctx *gin.Context) {
	var allPhotos = []string{}
	users, err := server.store.ListUserByAdmin(ctx)
	if err == nil && len(users) > 0 {
		for _, u := range users {
			allPhotos = append(allPhotos, u.Photo)
		}
	}
	res := ListFirebasePhoto{
		Photos: allPhotos,
	}
	ctx.JSON(http.StatusOK, res)
}

//func (server *Server) ListPhotoIdentityAdmin(ctx *gin.Context) {
//	var allPhotos = []string{}
//	ids, err := server.store.ListIdentityByAdmin(ctx)
//	if err == nil && len(ids) > 0 {
//		for _, id := range ids {
//			if id.IDPhoto != "none" {
//				allPhotos = append(allPhotos, id.IDPhoto)
//			}
//			if id.IDBackPhoto != "none" {
//				allPhotos = append(allPhotos, id.IDBackPhoto)
//			}
//			if id.FacialPhoto != "none" {
//				allPhotos = append(allPhotos, id.FacialPhoto)
//			}
//		}
//	}
//	res := ListFirebasePhoto{
//		Photos: allPhotos,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

func (server *Server) ListPhotoOptionAdmin(ctx *gin.Context) {
	var allPhotos = []string{}
	optionPhotos, err := server.store.ListOptionPhotoByAdmin(ctx)
	if err == nil && len(optionPhotos) > 0 {
		for _, op := range optionPhotos {
			if len(op.Photo) > 0 {
				allPhotos = append(allPhotos, op.Photo...)
			}
			allPhotos = append(allPhotos, op.CoverImage)
		}
	}
	res := ListFirebasePhoto{
		Photos: allPhotos,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListPhotoEventCheckInStepAdmin(ctx *gin.Context) {
	var allPhotos = []string{}
	eventCheckInSteps, err := server.store.ListEventCheckInStepByAdmin(ctx)
	if err == nil && len(eventCheckInSteps) > 0 {
		for _, e := range eventCheckInSteps {
			allPhotos = append(allPhotos, e.Photo)
		}
	}

	res := ListFirebasePhoto{
		Photos: allPhotos,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListPhotoCheckInStepAdmin(ctx *gin.Context) {
	var allPhotos = []string{}
	checkInSteps, err := server.store.ListCheckInStepByAdmin(ctx)
	if err == nil && len(checkInSteps) > 0 {
		for _, optionStep := range checkInSteps {
			allPhotos = append(allPhotos, optionStep.Photo)
		}
	}

	res := ListFirebasePhoto{
		Photos: allPhotos,
	}
	ctx.JSON(http.StatusOK, res)
}

// Update

func (server *Server) UpdatePhotoUserAdmin(ctx *gin.Context) {
	var req UpdateFirebasePhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	users, err := server.store.ListUserByAdmin(ctx)
	if err == nil && len(users) > 0 {
		for _, u := range users {
			if u.Photo == req.Actual {
				image := fmt.Sprintf("%v*%v", req.Actual, req.Update)
				_, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
					Image: pgtype.Text{
						String: image,
						Valid:  true,
					},
					ID: u.ID,
				})
				if err != nil {
					log.Printf("Error at UpdatePhotoUserAdmin %v\n", err.Error())
				}
				break
			}
		}
	}
	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

//func (server *Server) UpdatePhotoIdentityAdmin(ctx *gin.Context) {
//	var req UpdateFirebasePhoto
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("error at UpdateIdentity in ShouldBindJSON: %v \n", err)
//		err = fmt.Errorf("there was an error while processing your inputs please try again later")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	ids, err := server.store.ListIdentityByAdmin(ctx)
//	if err == nil && len(ids) > 0 {
//		for _, id := range ids {
//			if id.IDPhoto != "none" {
//				allPhotos = append(allPhotos, id.IDPhoto)
//			}
//			if id.IDBackPhoto != "none" {
//				allPhotos = append(allPhotos, id.IDBackPhoto)
//			}
//			if id.FacialPhoto != "none" {
//				allPhotos = append(allPhotos, id.FacialPhoto)
//			}
//		}
//	}
//	res := ListFirebasePhoto{
//		Photos: allPhotos,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

func (server *Server) UpdatePhotoOptionAdmin(ctx *gin.Context) {
	var req UpdateFirebasePhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateOption in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionPhotos, err := server.store.ListOptionPhotoByAdmin(ctx)
	if err == nil && len(optionPhotos) > 0 {
		for _, op := range optionPhotos {
			if op.CoverImage == req.Actual {
				image := fmt.Sprintf("%v*%v", req.Actual, req.Update)
				_, err = server.store.UpdateOptionInfoMainImage(ctx, db.UpdateOptionInfoMainImageParams{
					OptionID:  op.OptionID,
					MainImage: image,
				})
				if err != nil {
					log.Printf("Error at UpdatePhotoOptionAdmin cover %v\n", err.Error())
				}
				break
			}
			for _, basicPhotos := range op.Photo {
				if basicPhotos == req.Actual {
					image := fmt.Sprintf("%v*%v", req.Actual, req.Update)
					if !tools.IsInList(op.Images, image) {
						newPhoto := []string{req.Update}
						newPhoto = append(newPhoto, op.Images...)
						_, err = server.store.UpdateOptionInfoImages(ctx, db.UpdateOptionInfoImagesParams{
							OptionID:    op.OptionID,
							Images: newPhoto,
						})
						if err != nil {
							log.Printf("Error at UpdatePhotoOptionAdmin photo %v\n", err.Error())
						}
						break
					}
				}
			}
		}
	}
	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdatePhotoEventCheckInStepAdmin(ctx *gin.Context) {
	var req UpdateFirebasePhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateEventCheckInStep in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventCheckInSteps, err := server.store.ListEventCheckInStepByAdmin(ctx)
	if err == nil && len(eventCheckInSteps) > 0 {
		for _, e := range eventCheckInSteps {
			if e.Photo == req.Actual && e.PublicPhoto != req.Update {
				_, err = server.store.UpdateEventCheckInStepPublicPhoto(ctx, db.UpdateEventCheckInStepPublicPhotoParams{
					ID:          e.ID,
					PublicPhoto: req.Update,
				})
				if err != nil {
					log.Printf("Error at UpdatePhotoOptionAdmin photo %v\n", err.Error())
				}
				break
			}
		}
	}

	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdatePhotoCheckInStepAdmin(ctx *gin.Context) {
	var req UpdateFirebasePhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateCheckInStep in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	checkInSteps, err := server.store.ListCheckInStepByAdmin(ctx)
	if err == nil && len(checkInSteps) > 0 {
		for _, optionStep := range checkInSteps {
			if optionStep.Photo == req.Actual && optionStep.PublicPhoto != req.Update {
				_, err = server.store.UpdateCheckInStepPublicPhoto(ctx, db.UpdateCheckInStepPublicPhotoParams{
					ID:          optionStep.ID,
					PublicPhoto: req.Update,
				})
				if err != nil {
					log.Printf("Error at UpdatePhotoOptionAdmin photo %v\n", err.Error())
				}
				break
			}
		}
	}

	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}
