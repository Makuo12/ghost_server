package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func HandleImageMetaData(ctx context.Context, server *Server) func() {
	return func() {
		photos, err := server.store.ListAllUserPhotos(ctx)
		if err != nil || len(photos) == 0 {
			if err != nil {
				log.Println("HandleImageMetaData err: ", err)
			}
			return
		}
		for _, p := range photos {
			// First we handle the photo
			if tools.ServerStringEmpty(p.Photo) {
				log.Println("no object found here try again")

				continue
			}
			UpdateFireStorageMeta(ctx, server, p.Photo)

			if tools.ServerStringEmpty(p.FacialPhoto) {
				log.Println("no object found here try again")
				continue
			}
			UpdateFireStorageMeta(ctx, server, p.FacialPhoto)

			if tools.ServerStringEmpty(p.IDPhoto) {
				log.Println("no object found here try again")
				continue
			}
			UpdateFireStorageMeta(ctx, server, p.IDPhoto)
		}
	}

}

func UpdateFireStorageMeta(ctx context.Context, server *Server, object string) {
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	o := server.Bucket.Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		log.Printf("object.Attrs error: %v\n", err)
		return
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Create metadata
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		ContentType: "image/jpeg", // Set the content type, adjust accordingly

		//Metadata: map[string]string{
		//	"key1": "value1",
		//	"key2": "value2",
		//	// Add any custom metadata key-value pairs as needed
		//},
	}
	// Update the metadata for the file
	if _, err := o.Update(contextOne, objectAttrsToUpdate); err != nil {
		log.Fatalf("Failed to update metadata: %v", err)
		return
	}
	fmt.Println("Metadata updated successfully")
}

func UpdateFireStoragePublicUrl(ctx context.Context, server *Server, coverImage string, photos []string, optionID uuid.UUID) {
	// Generate a signed URL with an expiration time (e.g., 1 hour)
	//url, err := fileRef.SignedURL(context.Background(), time.Now().Add(1*time.Hour), nil)
	//if err != nil {
	//	log.Fatalf("Failed to generate signed URL: %v", err)
	//}

	//// Extract the token from the URL
	//token := url.Query().Get("token")

	//// Construct the desired download URL
	//downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s.appspot.com/o/%s?alt=media&token=%s", projectID, filePath, token)

	//fmt.Printf("Download URL: %s\n", downloadURL)
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	o := server.Bucket.Object(coverImage)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		err = fmt.Errorf("coverImage.Attrs: %v", err)
		return
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
	// Create metadata
	// Get the public URL
	publicURL, err := o.Attrs(ctx)
	if err != nil {
		log.Fatalf("Failed to get file attributes: %v", err)
		return
	}
	_, err = server.store.UpdateOptionInfoPhotoCoverUrl(ctx, db.UpdateOptionInfoPhotoCoverUrlParams{
		OptionID:         optionID,
		PublicCoverImage: publicURL.MediaLink,
	})
	if err != nil {
		log.Fatalf("UpdateOptionInfoPhotoCoverUrl: %v", err)
		return
	}
	// We handle normal photos
	var publicURLs []string
	for _, p := range photos {
		o := server.Bucket.Object(p)
		attrs, err := o.Attrs(ctx)
		if err != nil {
			err = fmt.Errorf("coverImage.Attrs: %v", err)
			continue
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
		// Create metadata
		// Get the public URL
		publicURL, err := o.Attrs(contextOne)
		if err != nil {
			log.Fatalf("Failed to get file attributes: %v", err)
			continue
		}
		publicURLs = append(publicURLs, publicURL.Name)
	}
	if len(publicURLs) != 0 {
		_, err = server.store.UpdateOptionInfoPhotoOnlyUrl(ctx, db.UpdateOptionInfoPhotoOnlyUrlParams{
			OptionID:    optionID,
			PublicPhoto: publicURLs,
		})
		if err != nil {
			log.Fatalf("UpdateOptionInfoPhotoOnlyUrl: %v", err)
			return
		}
	}
}

func constructDownloadURL(projectID, filePath, token string) string {
	baseURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s.appspot.com/o/%s", projectID, url.PathEscape(filePath))

	queryParams := url.Values{}
	queryParams.Set("alt", "media")
	queryParams.Set("token", token)

	fullURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Error parsing base URL: %v", err)
	}

	fullURL.RawQuery = queryParams.Encode()
	return fullURL.String()
}

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
				_, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
					PublicPhoto: pgtype.Text{
						String: req.Update,
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
				_, err = server.store.UpdateOptionInfoPhotoCoverUrl(ctx, db.UpdateOptionInfoPhotoCoverUrlParams{
					OptionID:         op.OptionID,
					PublicCoverImage: req.Update,
				})
				if err != nil {
					log.Printf("Error at UpdatePhotoOptionAdmin cover %v\n", err.Error())
				}
				break
			}
			for _, basicPhotos := range op.Photo {
				if basicPhotos == req.Actual {
					if !tools.IsInList(op.PublicPhoto, req.Update) {
						newPhoto := []string{req.Update}
						newPhoto = append(newPhoto, op.PublicPhoto...)
						_, err = server.store.UpdateOptionInfoPhotoOnlyUrl(ctx, db.UpdateOptionInfoPhotoOnlyUrlParams{
							OptionID:    op.OptionID,
							PublicPhoto: newPhoto,
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
