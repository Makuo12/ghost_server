package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleReserveOptionHost(selection string, server *Server, ctx *gin.Context, user db.User) (res []ReserveHostItem, hasData bool, err error) {
	data, err := server.store.ListChargeOptionReferenceHost(ctx, db.ListChargeOptionReferenceHostParams{
		CoUserID: tools.UuidToString(user.UserID),
		HostID:   user.ID,
	})
	if err != nil {
		log.Printf("Error at  HandleReserveCoHost in ListOptionHostInfo err: %v, user: %v\n", err, user.ID)
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
		}
		return
	}

	for _, d := range data {
		var canScanCode bool
		var canReserve bool
		var dOptionID string
		currentTime := time.Now()
		currentDateOnlyString := tools.ConvertTimeToStringDateOnly(currentTime)
		startDate := tools.ConvertDateOnlyToString(d.StartDate)
		endDate := tools.ConvertDateOnlyToString(d.EndDate)
		var leaveBeforeString string = "12:00"
		if !tools.ServerStringEmpty(d.LeaveBefore) {
			leaveBeforeString = d.LeaveBefore
		}
		leaveBefore, err := tools.ConvertDateTimeStringToTime(endDate, leaveBeforeString)
		if err != nil {
			log.Printf("Error at  HandleTime in ConvertDateTimeStringToTim err: %v, user: %v, startTimeType: %v\n", err, user.ID, "leave_before")
			continue
		}
		if d.HostType == "main_host" {
			dOptionID = tools.UuidToString(d.OptionID)
		} else {
			dOptionID = tools.UuidToString(d.CoHostID)
		}
		if d.HostType == "main_host" || d.ScanCode {
			canScanCode = true
		}
		if d.HostType == "main_host" || d.Reservations {
			canReserve = true
		}
		if !tools.ServerStringEmpty(d.ArriveAfter) {
			hasItem, dataRes, err := HandleTime(selection, user, d, startDate, endDate, d.ArriveAfter, currentTime, "arrive_after", leaveBefore, currentDateOnlyString, canReserve, canScanCode, dOptionID)
			if err == nil && hasItem {
				res = append(res, dataRes)
				continue
			}
		} else if !tools.ServerStringEmpty(d.ArriveBefore) {
			hasItem, dataRes, err := HandleTime(selection, user, d, startDate, endDate, d.ArriveBefore, currentTime, "arrive_before", leaveBefore, currentDateOnlyString, canReserve, canScanCode, dOptionID)
			if err == nil && hasItem {
				res = append(res, dataRes)
				continue
			}
		} else {
			hasItem, dataRes, err := HandleTime(selection, user, d, startDate, endDate, "08:00", currentTime, "none", leaveBefore, currentDateOnlyString, canReserve, canScanCode, dOptionID)
			if err == nil && hasItem {
				res = append(res, dataRes)
				continue
			}
		}
	}
	if len(res) > 0 {
		hasData = true
	}
	return
}

func HandleTime(selection string, user db.User, d db.ListChargeOptionReferenceHostRow, startDate string, endDate string, timeString string, currentTime time.Time, startTimeType string, leaveBefore time.Time, currentDateOnlyString string, canReserve bool, canScanCode bool, dOptionID string) (hasItem bool, data ReserveHostItem, err error) {
	// We want to add 2 hours to the time
	timeData, err := tools.ConvertDateTimeStringToTime(startDate, timeString)
	if err != nil {
		log.Printf("Error at  HandleTime in ConvertDateTimeStringToTime err: %v, user: %v, startTimeType: %v\n", err, user.ID, startTimeType)
	} else {

		switch selection {
		case "currently_hosting":
			timeData = tools.AddHoursToTime(timeData, 2)
			if startDate != currentDateOnlyString && currentTime.After(timeData) && currentTime.Before(leaveBefore) {
				data = ReserveHostItem{
					UserID:         tools.UuidToString(d.UserID),
					StartDate:      startDate,
					EndDate:        endDate,
					OptionID:       dOptionID,
					FirstName:      d.FirstName,
					HostNameOption: d.HostNameOption,
					UserPhoto:      d.Photo,
					ArriveAfter:    d.ArriveAfter,
					ArriveBefore:   d.ArriveBefore,
					LeaveBefore:    d.LeaveBefore,
					TimeZone:       d.TimeZone,
					StartTimeType:  startTimeType,
					HostMethod:     d.HostType,
					ReferenceID:    tools.UuidToString(d.ReferenceID),
					CanScanCode:    canScanCode,
					CanReserve:     canReserve,
					OptionStatus:   d.Status,
					CoverImage:     d.CoverImage,
				}
				hasItem = true

			} else {
				hasItem = false
			}
			return
		case "coming_soon":
			if (startDate == currentDateOnlyString) && currentTime.Before(leaveBefore) {
				data = ReserveHostItem{
					UserID:         tools.UuidToString(d.UserID),
					StartDate:      startDate,
					EndDate:        endDate,
					OptionID:       dOptionID,
					FirstName:      d.FirstName,
					HostNameOption: d.HostNameOption,
					UserPhoto:      d.Photo,
					ArriveAfter:    d.ArriveAfter,
					ArriveBefore:   d.ArriveBefore,
					LeaveBefore:    d.LeaveBefore,
					TimeZone:       d.TimeZone,
					StartTimeType:  startTimeType,
					HostMethod:     d.HostType,
					ReferenceID:    tools.UuidToString(d.ReferenceID),
					CanScanCode:    canScanCode,
					CanReserve:     canReserve,
					OptionStatus:   d.Status,
					CoverImage:     d.CoverImage,
				}
				hasItem = true

			} else {
				hasItem = false
			}
			return
		case "upcoming":
			if startDate != currentDateOnlyString && currentTime.Before(timeData) && currentTime.Before(leaveBefore) {
				data = ReserveHostItem{
					UserID:         tools.UuidToString(d.UserID),
					StartDate:      startDate,
					EndDate:        endDate,
					OptionID:       dOptionID,
					FirstName:      d.FirstName,
					HostNameOption: d.HostNameOption,
					UserPhoto:      d.Photo,
					ArriveAfter:    d.ArriveAfter,
					ArriveBefore:   d.ArriveBefore,
					LeaveBefore:    d.LeaveBefore,
					TimeZone:       d.TimeZone,
					StartTimeType:  startTimeType,
					HostMethod:     d.HostType,
					ReferenceID:    tools.UuidToString(d.ReferenceID),
					CanScanCode:    canScanCode,
					CanReserve:     canReserve,
					OptionStatus:   d.Status,
					CoverImage:     d.CoverImage,
				}
				hasItem = true

			} else {
				hasItem = false
			}
		case "checking_out":
			if endDate == currentDateOnlyString && currentTime.After(timeData) && currentTime.After(leaveBefore) {
				data = ReserveHostItem{
					UserID:         tools.UuidToString(d.UserID),
					StartDate:      startDate,
					EndDate:        endDate,
					OptionID:       dOptionID,
					FirstName:      d.FirstName,
					HostNameOption: d.HostNameOption,
					UserPhoto:      d.Photo,
					ArriveAfter:    d.ArriveAfter,
					ArriveBefore:   d.ArriveBefore,
					LeaveBefore:    d.LeaveBefore,
					TimeZone:       d.TimeZone,
					StartTimeType:  startTimeType,
					HostMethod:     d.HostType,
					ReferenceID:    tools.UuidToString(d.ReferenceID),
					CanScanCode:    canScanCode,
					CanReserve:     canReserve,
					OptionStatus:   d.Status,
					CoverImage:     d.CoverImage,
				}
				hasItem = true

			} else {
				hasItem = false
			}
			return
		case "checked_out":
			if endDate != currentDateOnlyString && currentTime.After(timeData) && currentTime.After(leaveBefore) {
				data = ReserveHostItem{
					UserID:         tools.UuidToString(d.UserID),
					StartDate:      startDate,
					EndDate:        endDate,
					OptionID:       dOptionID,
					FirstName:      d.FirstName,
					HostNameOption: d.HostNameOption,
					UserPhoto:      d.Photo,
					ArriveAfter:    d.ArriveAfter,
					ArriveBefore:   d.ArriveBefore,
					LeaveBefore:    d.LeaveBefore,
					TimeZone:       d.TimeZone,
					StartTimeType:  startTimeType,
					HostMethod:     d.HostType,
					ReferenceID:    tools.UuidToString(d.ReferenceID),
					CanScanCode:    canScanCode,
					CanReserve:     canReserve,
					OptionStatus:   d.Status,
					CoverImage:     d.CoverImage,
				}
				hasItem = true

			} else {
				hasItem = false
			}
			return
		}
	}
	return
}

func ListReservationDetailResOffset(data []ReserveHostItem, offset int, limit int) []ReserveHostItem {
    // If offset is greater than or equal to the length of data, return an empty slice.
    if offset >= len(data) {
        return []ReserveHostItem{}
    }

    // Calculate the end index based on offset and limit.
    end := offset + limit

    // If the end index is greater than the length of data, set it to the length of data.
    if end > len(data) {
        end = len(data)
    }

    // Return a subset of data starting from the offset and up to the end index.
    return data[offset:end]
}


