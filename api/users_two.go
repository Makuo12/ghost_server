package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/val"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nyaruka/phonenumbers"
)

func (server *Server) UpdateIdentity(ctx *gin.Context) {
	var req UpdateIdentityParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateIdentity in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(req.IDPhoto) < 1 || len(req.FacialPhoto) < 1 {
		err = fmt.Errorf("photo is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	currentIdentity, err := server.store.GetIdentity(ctx, user.ID)
	if err != nil {
		log.Printf("Error at  UpdateIdentity in GetIdentity err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not update your identity profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	idPhotoList := tools.ServerListToDB(currentIdentity.IDPhotoList)
	facialPhotoList := tools.ServerListToDB(currentIdentity.FacialPhotoList)
	idBackPhotoList := tools.ServerListToDB(currentIdentity.IDBackPhotoList)
	idPhotoList = append(idPhotoList, req.IDPhoto)
	facialPhotoList = append(facialPhotoList, req.FacialPhoto)
	if !tools.ServerStringEmpty(req.IDBackPhoto) {
		idBackPhotoList = append(idBackPhotoList, req.IDBackPhoto)
	} else {
		if len(idBackPhotoList) == 0 {
			idBackPhotoList = append(idBackPhotoList, "none")
		}
	}
	// We just want to update the identity
	identity, err := server.store.UpdateIdentity(ctx, db.UpdateIdentityParams{
		Country:         req.Country,
		Type:            req.Type,
		IDPhoto:         req.IDPhoto,
		FacialPhoto:     req.FacialPhoto,
		Status:          "processing",
		IDPhotoList:     idPhotoList,
		FacialPhotoList: facialPhotoList,
		IDBackPhotoList: idBackPhotoList,
		IsVerified:      false,
		UserID:          user.ID,
	})
	if err != nil {
		log.Printf("Error at  UpdateIdentity in UpdateIdentity err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not update your identity profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateIdentityRes{
		Status:   identity.Status,
		Verified: identity.IsVerified,
	}
	log.Printf("UpdateIdentity successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateEmContact(ctx *gin.Context) {
	var req CreateEmContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateEmContact in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !tools.ServerStringEmpty(req.Email) {
		if !val.ValidateEmail(req.Email) {
			err := fmt.Errorf("email not in the right format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if tools.ServerStringEmpty(req.Relationship) {
		err := fmt.Errorf("your relationship to the person cannot be empty")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(req.Code) == 0 || len(req.PhoneNumber) == 0 || len(req.DialCountry) == 0 || req.DialCountry == "none" || req.PhoneNumber == "none" || req.Code == "none" {
		err := fmt.Errorf("phone number format wrong, please try entering your phone number again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	num, err := phonenumbers.Parse(req.PhoneNumber, req.Code)
	if err != nil {
		log.Printf("Error at VerifyPhoneNumber in phonenumbers: %v \n", err.Error())
		err = fmt.Errorf("there was an error while checking your number, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	number := fmt.Sprintf("%d_%d", num.GetCountryCode(), num.GetNationalNumber())
	dial_code := strconv.FormatInt(int64(num.GetCountryCode()), 10) // 234
	dial_country := req.DialCountry                                 // Nigeria
	// We just want to update the em contact
	_, err = server.store.GetEmContactByPhone(ctx, number)
	if err != nil {
		// We expect this error
		log.Printf("Error at  CreateEmContact in GetEmContactByPhone err: %v, user: %v\n", err, user.ID)
	} else {
		// We send an error because it means the data exist
		err = fmt.Errorf("this contact already exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	emContact, err := server.store.CreateEmContact(ctx, db.CreateEmContactParams{
		Name:         req.Name,
		Relationship: strings.ToLower(req.Relationship),
		Email:        req.Email,
		UserID:       user.ID,
		DialCode:     dial_code,
		DialCountry:  dial_country,
		PhoneNumber:  number,
		Language:     req.Language,
	})
	if err != nil {
		log.Printf("Error at  CreateEmContact in CreateEmContact err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not update your identity profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := EmContactDetail{
		Name: emContact.Name,
		ID:   tools.UuidToString(emContact.ID),
	}
	log.Printf("CreateEmContact successfully (%v) id: %v\n em_name: %v", user.Email, user.ID, res.Name)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveEmContact(ctx *gin.Context) {
	var req RemoveEmContactParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at RemoveEmContact in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	emID, err := tools.StringToUuid(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.RemoveEmContact(ctx, emID)
	if err != nil {
		log.Printf("Error at  RemoveEmContact in RemoveEmContact err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not update your identity profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := RemoveResponse{
		Success: true,
	}
	log.Printf("RemoveEmContact successfully (%v) id: %v\n em_name: %v", user.Email, user.ID, emID)
	ctx.JSON(http.StatusOK, res)
}

// Get the user email and firebase password
func (server *Server) GetFireEmailAndPassword(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetFireEmailAndPasswordRes{
		Email:        user.Email,
		FireFight:    user.FirebasePassword,
		ProfilePhoto: user.Photo,
	}
	log.Printf("GetFireEmailAndPassword (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Get the user email and firebase password
func (server *Server) GetProfileUser(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	var contacts []EmContactDetail
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	gid, err := server.store.GetIdentity(ctx, user.ID)
	if err != nil {
		log.Printf("Error at  ProfileUser in GetIdentity err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not process your identity profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	emContacts, err := server.store.ListEmContact(ctx, user.ID)
	if err != nil {
		log.Printf("Error at  ProfileUser in ListEmContact err: %v, user: %v\n", err, user.ID)
		contacts = append(contacts, EmContactDetail{Name: "", ID: ""})
	} else {
		if len(emContacts) == 0 {
			contacts = append(contacts, EmContactDetail{Name: "", ID: ""})
		} else {
			for _, c := range emContacts {
				contacts = append(contacts, EmContactDetail{Name: c.Name, ID: tools.UuidToString(c.ID)})
			}
		}
	}
	var phoneNumber string
	if tools.ServerStringEmpty(user.PhoneNumber) {
		phoneNumber = ""
	} else {
		// Phone number use _
		phoneData := strings.Split(user.PhoneNumber, "_")
		phoneNumber = fmt.Sprintf("%v%v", phoneData[0], tools.HandlePhoneNumber(phoneData[1]))
	}
	email := tools.HandleEmail(user.Email)
	res := ProfileUserRes{
		Email:       email,
		PhoneNumber: phoneNumber,
		Status:      gid.Status,
		Verified:    gid.IsVerified,
		PhoneCode:   user.DialCode,
		DateOfBirth: tools.ConvertDateOnlyToString(user.DateOfBirth),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		EmContacts:  contacts,
		Currency:    user.Currency,
		UserTwoID:   tools.UuidToString(user.UserID),
	}
	log.Printf("ProfileUserRes sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateUserProfile(ctx *gin.Context) {
	var req UserProfileParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateUserProfile in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var userProfile db.UsersProfile
	switch req.Type {
	case "work":
		userProfile, err = server.store.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
			Work: pgtype.Text{
				String: req.Work,
				Valid:  true,
			},
			UserID: user.ID,
		})
		if err != nil {
			log.Printf("Error at  UpdateUserProfile in Work  in UpdateUserProfile err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not update your work info")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "bio":
		userProfile, err = server.store.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
			Bio: pgtype.Text{
				String: req.Bio,
				Valid:  true,
			},
			UserID: user.ID,
		})
		if err != nil {
			log.Printf("Error at  UpdateUserProfile in Bio in UpdateUserProfile err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not update your bio info")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "language":
		var langs = tools.HandleDBList(req.Languages)

		userProfile, err = server.store.UpdateUserProfileLang(ctx, db.UpdateUserProfileLangParams{
			Languages: langs,
			UserID:    user.ID,
		})
		if err != nil {
			log.Printf("Error at  UpdateUserProfile in UpdateUserProfileLang err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not update your languages info")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	res := UserProfileParams{
		Work:      userProfile.Work,
		Bio:       userProfile.Bio,
		Languages: userProfile.Languages,
	}
	log.Printf("UpdateUserProfile successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateUserProfilePhoto(ctx *gin.Context) {
	var req UpdateProfilePhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateUserProfilePhoto in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !tools.ServerStringEmpty(user.Photo) {
		err = RemoveFirebasePhoto(server, ctx, user.Photo)
		if err != nil {
			log.Printf("Error at UpdateUserProfilePhoto in RemoveFirebasePhoto err: %v, user: %v\n", err, user.ID)
		}
	}
	userUpdate, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		Photo: pgtype.Text{
			String: req.ProfilePhoto,
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("Error at UpdateUserProfilePhoto in UpdateUser err: %v, user: %v photoID: %v\n", err, user.ID, req.ProfilePhoto)
		err = fmt.Errorf("could not update your profile photo")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UpdateProfilePhotoParams{
		ProfilePhoto: userUpdate.Photo,
	}
	log.Printf("UpdateUserProfilePhoto successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Get the user currency
func (server *Server) GetUserCurrency(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetUserCurrencyRes{
		Currency: user.Currency,
	}
	log.Printf("GetUserCurrency sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Get the user profilePhoto
func (server *Server) GetUserProfilePhoto(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetUserProfilePhotoRes{
		ProfilePhoto: user.Photo,
	}
	log.Printf("GetUserProfilePhoto sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetUserIsHost(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	unreadMsgs := UserUnreadMessages(ctx, server, user, "GetUserIsHost")
	unreadNotifications := UserUnreadNotifications(ctx, server, user, "GetUserIsHost")
	res := GetUserIsHostRes{
		IsHost:              userIsHost,
		HasIncomplete:       hasIncomplete,
		UnreadMessages:      unreadMsgs,
		UnreadNotifications: int(unreadNotifications),
		ProfileImage:        user.Photo,
	}
	log.Printf("GetUserIsHost sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateUserLocation(ctx *gin.Context) {
	var req CreateUpdateUserLocationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateUpdateUserLocation in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res CreateUpdateUserLocationRes
	_, err = server.store.GetUserLocationHalf(ctx, user.ID)
	if err != nil {
		log.Printf("Error at CreateUpdateUserLocation in GetUserLocationHalf err: %v, user: %v", err, user.ID)
		// If there was an error we then create it
		lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
		lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
		geolocation := pgtype.Point{
			P:     pgtype.Vec2{X: lng, Y: lat},
			Valid: true,
		}
		location, err := server.store.CreateUserLocation(ctx, db.CreateUserLocationParams{
			UserID:      user.ID,
			Street:      req.Street,
			City:        req.City,
			State:       req.State,
			Country:     req.Country,
			Postcode:    req.Postcode,
			Geolocation: geolocation,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateUserLocation in CreateUserLocation err: %v, user: %v", err, user.ID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CreateUpdateUserLocationRes{
				State:   location.State,
				Country: location.Country,
			}
		}
	} else {
		lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
		lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
		geolocation := pgtype.Point{
			P:     pgtype.Vec2{X: lng, Y: lat},
			Valid: true,
		}
		// Not error then we want to create
		location, err := server.store.UpdateUserLocation(ctx, db.UpdateUserLocationParams{
			UserID:      user.ID,
			Street:      req.Street,
			City:        req.City,
			State:       req.State,
			Country:     req.Country,
			Postcode:    req.Postcode,
			Geolocation: geolocation,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateUserLocation in UpdateUserLocation err: %v, user: %v", err, user.ID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CreateUpdateUserLocationRes{
				State:   location.State,
				Country: location.Country,
			}
		}
	}

	log.Printf("CreateUpdateUserLocation successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveUserLocation(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}
	err = server.store.RemoveUserLocation(ctx, user.ID)
	if err != nil {
		log.Printf("Error at RemoveUserLocation in RemoveUserLocation err: %v, user: %v", err, user.ID)
		err := fmt.Errorf("was unable to remove your location, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := RemoveResponse{
		Success: true,
	}
	log.Printf("RemoveUserLocation successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetUserProfileDetail(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var emailConfirmed bool
	var phoneConfirmed bool
	var languages []string
	var options []ProfileDetailOption
	var yearJoined string
	var state string
	var country string
	customOption := ProfileDetailOption{
		OptionUserID: "",
		Name:         "",
		CoverImage:   "",
		Type:         "",
		MainOption:   "",
	}
	profile, err := server.store.GetUserProfile(ctx, user.ID)
	if err != nil {
		log.Printf("Error at GetUserProfileDetail in GetUserProfile err: %v, user: %v", err, user.ID)
		err := fmt.Errorf("was unable to fetch your profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(profile.Languages) == 0 {
		languages = []string{"none"}
	} else {
		languages = profile.Languages
	}
	identity, err := server.store.GetIdentity(ctx, user.ID)
	if err != nil {
		log.Printf("Error at GetUserProfileDetail in GetIdentity err: %v, user: %v", err, user.ID)
		err := fmt.Errorf("was unable to fetch your identity")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	location, err := server.store.GetUserLocationHalf(ctx, user.ID)
	if err != nil {
		log.Printf("Error at GetUserProfileDetail in GetUserLocationHalf err: %v, user: %v", err, user.ID)
		state = ""
		country = ""
	} else {
		state = location.State
		country = location.Country
	}
	yearJoined = tools.ConvertTimeToYear(user.CreatedAt)

	// Confirm email
	if !tools.ServerStringEmpty(user.Email) {
		emailConfirmed = true
	}
	// Confirm phone number
	if !tools.ServerStringEmpty(user.PhoneNumber) && !tools.ServerStringEmpty(user.DialCode) {
		phoneConfirmed = true
	}
	optionsData, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil || len(optionsData) == 0 {
		if err != nil {
			log.Printf("Error at GetUserProfileDetail in ListOptionInfoNoLimit err: %v, user: %v", err, user.ID)
		}
		options = append(options, customOption)
	} else {
		for _, option := range optionsData {
			var optionType string = ""
			switch option.MainOptionType {
			case "options":
				optionType = HandleSqlNullString(option.TypeOfShortlet)
			case "events":
				optionType = HandleSqlNullString(option.EventType)
			}
			data := ProfileDetailOption{
				OptionUserID: tools.UuidToString(option.OptionUserID),
				Name:         option.HostNameOption,
				CoverImage:   HandleSqlNullString(option.CoverImage),
				Type:         optionType,
				MainOption:   option.MainOptionType,
			}
			options = append(options, data)
		}
	}
	res := GetUserProfileDetailRes{
		Status:         identity.Status,
		Verified:       identity.IsVerified,
		EmailConfirmed: emailConfirmed,
		PhoneConfirmed: phoneConfirmed,
		Bio:            profile.Bio,
		Work:           profile.Work,
		Languages:      languages,
		State:          state,
		Country:        country,
		YearJoined:     yearJoined,
		UserTwoId:      tools.UuidToString(user.UserID),
		Options:        options,
	}
	log.Printf("GetUserProfileDetail successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Only a user that is logged in that can send a feedback
func (server *Server) CreateFeedback(ctx *gin.Context) {
	var req CreateFeedbackParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateFeedback in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.CreateFeedback(ctx, db.CreateFeedbackParams{
		UserID:     user.ID,
		Subject:    req.Subject,
		SubSubject: "none",
		Detail:     req.Detail,
	})
	if err != nil {
		log.Printf("Error at CreateFeedback in CreateFeedback err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("your request was unsuccessful please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UserResponseMsg{
		Success: true,
	}
	log.Printf("CreateFeedback successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// CreateHelpUser this is for user that are logged in
func (server *Server) CreateHelpUser(ctx *gin.Context) {
	var req CreateHelpUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateHelpUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.CreateHelp(ctx, db.CreateHelpParams{
		Email:      user.Email,
		Subject:    req.Subject,
		SubSubject: "none",
		Detail:     req.Detail,
	})
	if err != nil {
		log.Printf("Error at CreateHelpUser in CreateHelpUser err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("your request was unsuccessful please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UserResponseMsg{
		Success: true,
	}
	log.Printf("CreateHelpUser successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Help for users that don't have permissions
func (server *Server) CreateHelp(ctx *gin.Context) {
	var req CreateHelpParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateHelp in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.CreateHelp(ctx, db.CreateHelpParams{
		Email:      req.Email,
		Subject:    req.Subject,
		SubSubject: "none",
		Detail:     req.Detail,
	})
	if err != nil {
		log.Printf("Error at CreateHelp in CreateHelp err: %v, user: %v\n", err, req.Email)
		err = fmt.Errorf("your request was unsuccessful please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UserResponseMsg{
		Success: true,
	}
	log.Printf("CreateHelp successfully (%v) \n", req.Email)
	ctx.JSON(http.StatusOK, res)
}

// UpdateCurrencyUser
func (server *Server) UpdateCurrencyUser(ctx *gin.Context) {
	var req UpdateCurrencyUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateCurrencyUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
		Currency: pgtype.Text{
			String: req.Currency,
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("Error at UpdateCurrencyUser in UpdateUser err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("your request was unsuccessful please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := UpdateCurrencyUserParams{
		Currency: user.Currency,
	}
	log.Printf("UpdateCurrencyUser successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}
