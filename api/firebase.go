package api

import (
	"context"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

func HandleImageMetaData(ctx context.Context, server *Server) func() {
	return func() {
		photos, err := server.store.ListAllPhoto(ctx)
		if err != nil || len(photos) == 0 {
			if err != nil {
				log.Println("HandleImageMetaData err: ", err)
			}
			return
		}
		for _, p := range photos {
			// First we handle the cover_photo

			if tools.ServerStringEmpty(p.CoverImage) {
				err = fmt.Errorf("no object found here try again")
				continue
			}
			UpdateFireStorageMeta(ctx, server, p.CoverImage)
			for _, v := range p.Photo {
				UpdateFireStorageMeta(ctx, server, v)
			}
			UpdateFireStoragePublicUrl(ctx, server, p.CoverImage, p.Photo, p.OptionID)
		}
	}

}

func UpdateFireStorageMeta(ctx context.Context, server *Server, object string) {
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()
	o := server.Bucket.Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		err = fmt.Errorf("object.Attrs: %v", err)
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
	return
}

func UpdateFireStoragePublicUrl(ctx context.Context, server *Server, coverImage string, photos []string, optionID uuid.UUID) {
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
	publicURL, err := o.Attrs(contextOne)
	if err != nil {
		log.Fatalf("Failed to get file attributes: %v", err)
		return
	}
	_, err = server.store.UpdateOptionInfoPhotoCoverUrl(ctx, db.UpdateOptionInfoPhotoCoverUrlParams{
		OptionID:         optionID,
		PublicCoverImage: publicURL.Name,
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
