package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

func RemoveFirebasePhoto(server *Server, ctx context.Context, object string) (err error) {
	// First we delete cover photo
	if object == "none" || len(object) < 1 {
		err = fmt.Errorf("no object found here try again")
		return
	}
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	o := server.Bucket.Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		err = fmt.Errorf("object.Attrs: %v", err)
		return
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err = o.Delete(contextOne); err != nil {
		err = fmt.Errorf("Object(%q).Delete: %v", object, err)
		return
	}
	log.Printf("Object %v was deleted", object)
	return nil
}

//func UpdateFireStorageMeta(ctx context.Context, server *Server, object string) {
//	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
//	defer cancel()
//	o := server.Bucket.Object(object)
//	attrs, err := o.Attrs(ctx)
//	if err != nil {
//		log.Printf("object.Attrs error: %v\n", err)
//		return
//	}
//	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

//	// Create metadata
//	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
//		ContentType: "image/jpeg", // Set the content type, adjust accordingly

//		//Metadata: map[string]string{
//		//	"key1": "value1",
//		//	"key2": "value2",
//		//	// Add any custom metadata key-value pairs as needed
//		//},
//	}
//	// Update the metadata for the file
//	if _, err := o.Update(contextOne, objectAttrsToUpdate); err != nil {
//		log.Fatalf("Failed to update metadata: %v", err)
//		return
//	}
//	fmt.Println("Metadata updated successfully")
//}

//func UpdateFireStoragePublicUrl(ctx context.Context, server *Server, coverImage string, photos []string, optionID uuid.UUID) {
//	// Generate a signed URL with an expiration time (e.g., 1 hour)
//	//url, err := fileRef.SignedURL(context.Background(), time.Now().Add(1*time.Hour), nil)
//	//if err != nil {
//	//	log.Fatalf("Failed to generate signed URL: %v", err)
//	//}

//	//// Extract the token from the URL
//	//token := url.Query().Get("token")

//	//// Construct the desired download URL
//	//downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s.appspot.com/o/%s?alt=media&token=%s", projectID, filePath, token)

//	//fmt.Printf("Download URL: %s\n", downloadURL)
//	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
//	defer cancel()
//	o := server.Bucket.Object(coverImage)
//	attrs, err := o.Attrs(ctx)
//	if err != nil {
//		err = fmt.Errorf("coverImage.Attrs: %v", err)
//		return
//	}
//	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
//	// Create metadata
//	// Get the public URL
//	publicURL, err := o.Attrs(ctx)
//	if err != nil {
//		log.Fatalf("Failed to get file attributes: %v", err)
//		return
//	}
//	_, err = server.store.UpdateOptionInfoPhotoCoverUrl(ctx, db.UpdateOptionInfoPhotoCoverUrlParams{
//		OptionID:         optionID,
//		PublicCoverImage: publicURL.MediaLink,
//	})
//	if err != nil {
//		log.Fatalf("UpdateOptionInfoPhotoCoverUrl: %v", err)
//		return
//	}
//	// We handle normal photos
//	var publicURLs []string
//	for _, p := range photos {
//		o := server.Bucket.Object(p)
//		attrs, err := o.Attrs(ctx)
//		if err != nil {
//			err = fmt.Errorf("coverImage.Attrs: %v", err)
//			continue
//		}
//		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
//		// Create metadata
//		// Get the public URL
//		publicURL, err := o.Attrs(contextOne)
//		if err != nil {
//			log.Fatalf("Failed to get file attributes: %v", err)
//			continue
//		}
//		publicURLs = append(publicURLs, publicURL.Name)
//	}
//	if len(publicURLs) != 0 {
//		_, err = server.store.UpdateOptionInfoPhotoOnlyUrl(ctx, db.UpdateOptionInfoPhotoOnlyUrlParams{
//			OptionID:    optionID,
//			PublicPhoto: publicURLs,
//		})
//		if err != nil {
//			log.Fatalf("UpdateOptionInfoPhotoOnlyUrl: %v", err)
//			return
//		}
//	}
//}

//func constructDownloadURL(projectID, filePath, token string) string {
//	baseURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s.appspot.com/o/%s", projectID, url.PathEscape(filePath))

//	queryParams := url.Values{}
//	queryParams.Set("alt", "media")
//	queryParams.Set("token", token)

//	fullURL, err := url.Parse(baseURL)
//	if err != nil {
//		log.Fatalf("Error parsing base URL: %v", err)
//	}

//	fullURL.RawQuery = queryParams.Encode()
//	return fullURL.String()
//}
