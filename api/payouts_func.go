package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"time"
)

func GetDateAndAmount(data []db.ListChargeTicketReferencePayoutRow) (amount float64, date time.Time) {
	for i := 0; i < len(data); i++ {
		if i == 0 {
			date = data[i].TimePaid
		}
		amount += tools.ConvertStringToFloat(tools.IntToMoneyString(data[i].Amount))
		if date.Before(data[i].TimePaid) {
			date = data[i].TimePaid
		}

	}
	return
}
