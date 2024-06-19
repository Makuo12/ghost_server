package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type DeleteOptionResult struct {
	Success bool `json:"success"`
}

type DeleteOptionParams struct {
	OptionID     uuid.UUID   `json:"option_id"`
	OptionUserID uuid.UUID   `json:"option_user_id"`
	ChargeID     []uuid.UUID `json:"charge_id"`
}

func (store *SQLStore) DeleteOption(ctx context.Context, arg DeleteOptionParams, bucket *storage.BucketHandle) (DeleteOptionResult, error) {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		err = q.RemoveWishlistItemByOptionUserID(ctx, arg.OptionUserID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: wishlist item", err)
			return err
		}
		err = q.RemoveVid(ctx, arg.OptionUserID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: vid", err)
			return err
		}
		err = q.RemoveOptionInfoCategory(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option info category", err)
			return err
		}
		err = q.RemoveCompleteOptionInfoByID(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: complete option info", err)
			return err
		}
		err = q.RemoveOptionInfoDetails(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option info details", err)
			return err
		}
		err = q.RemoveOptionRemoveChargeByOptionID(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option charge", err)
			return err
		}

		err = q.RemoveOptionDiscountByOptionID(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option discount", err)
			return err
		}
		err = q.RemoveOptionPhotoCaptionByOptionID(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option caption", err)
			return err
		}
		err = q.RemoveAllOptionDateTime(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: option date time", err)
			return err
		}
		err = q.DeleteShortlet(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: shortlet", err)
			return err
		}
		err = q.RemoveCancelPolicy(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveCancelPolicy ", err)
			return err
		}
		err = q.RemoveOptionAvailabilitySetting(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionAvailabilitySetting ", err)
			return err
		}
		err = q.RemoveAllAmenity(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllAmenity ", err)
			return err
		}
		err = q.RemoveWifiDetail(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveWifiDetail ", err)
			return err
		}
		err = q.RemoveOptionPrice(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionPrice ", err)
			return err
		}
		err = q.RemoveLocation(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveLocation ", err)
			return err
		}
		err = q.RemoveAllSpaceAreas(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllSpaceAreas ", err)
			return err
		}
		err = q.RemoveAllThingToNote(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllThingToNote ", err)
			return err
		}
		err = q.RemoveOptionInfoStatus(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionInfoStatus ", err)
			return err
		}
		err = q.RemoveAllOptionMessage(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllOptionMessage ", err)
			return err
		}
		err = q.RemoveOptionQuestion(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionQuestion ", err)
			return err
		}
		err = q.RemoveOptionExtraInfo(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionExtraInfo ", err)
			return err
		}
		err = q.RemoveAllReportOption(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllReportOption ", err)
			return err
		}
		err = q.RemoveAllOptionRule(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllOptionRule ", err)
			return err
		}
		err = q.RemoveOptionBookMethod(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionBookMethod ", err)
			return err
		}
		err = q.RemoveBookRequirement(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveBookRequirement ", err)
			return err
		}
		err = q.RemoveAllOptionCOHostOptionID(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveAllOptionCOHostOptionID ", err)
			return err
		}
		err = q.RemoveCheckInOutDetail(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveCheckInOutDetail ", err)
			return err
		}
		err = q.RemoveOptionTripLength(ctx, arg.OptionID)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveOptionTripLength ", err)
			return err
		}
		for _, chargeID := range arg.ChargeID {
			err = HandleChargeOption(ctx, q, chargeID)
			if err != nil && err != ErrorRecordNotFound {
				return err
			}
		}
		err = HandleOptionPhoto(ctx, q, arg.OptionID, bucket)
		if err != nil && err != ErrorRecordNotFound {
			return err
		}
		return nil
	})
	if err != nil {
		return DeleteOptionResult{false}, err
	}
	return DeleteOptionResult{true}, err
}

func HandleChargeOption(ctx context.Context, q *Queries, chargeID uuid.UUID) error {
	err := q.RemoveChargeReview(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: charge review", err)
		return err
	}
	err = q.RemoveOptionReferenceInfo(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveOptionReferenceInfo", err)
		return err
	}
	err = q.RemoveMainRefunds(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveMainRefunds", err)
		return err
	}
	err = q.RemoveRefundPayout(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveRefundPayout", err)
		return err
	}
	err = q.RemoveRefund(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveRefund", err)
		return err
	}
	err = q.RemoveMainPayout(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveMainPayout", err)
		return err
	}
	err = q.RemoveRefund(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveRefund ", err)
		return err
	}
	err = q.RemoveRefund(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveRefund ", err)
		return err
	}
	err = q.RemoveChargeOptionReference(ctx, chargeID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveChargeOptionReference ", err)
		return err
	}
	return nil
}

func HandleOptionPhoto(ctx context.Context, q *Queries, optionID uuid.UUID, bucket *storage.BucketHandle) error {
	optionPhoto, err := q.GetOptionInfoPhoto(ctx, optionID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: GetOptionInfoPhoto", err)
		return err
	}
	photos := []string{}
	photos = append(photos, optionPhoto.Photo...)
	photos = append(photos, optionPhoto.CoverImage)
	// Delete all photos
	for _, p := range photos {
		err = RemoveFirebasePhoto(ctx, bucket, p)
		if err != nil && err != ErrorRecordNotFound {
			log.Println("err: RemoveFirebasePhoto", err)
		}
	}
	err = q.RemoveOptionInfoPhoto(ctx, optionID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveOptionInfoPhoto ", err)
		return err
	}
	step, err := q.GetCheckInStepByOptionID(ctx, optionID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: GetCheckInStepByOptionID", err)
		return err
	}

	err = RemoveFirebasePhoto(ctx, bucket, step.Photo)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveFirebasePhoto", err)
	}
	err = q.RemoveCheckInStepByOptionID(ctx, optionID)
	if err != nil && err != ErrorRecordNotFound {
		log.Println("err: RemoveCheckInStepByOptionID ", err)
		return err
	}
	return nil

}

func RemoveFirebasePhoto(ctx context.Context, bucket *storage.BucketHandle, object string) (err error) {
	// First we delete cover photo
	if object == "none" || len(object) < 1 {
		err = fmt.Errorf("no object found here try again")
		return
	}
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	o := bucket.Object(object)
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
