package api

import (
	"flex_server/val"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return val.IsSupportedCurrency(currency)
	}
	return false
}

var validEmail validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if email, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateEmail(email)
	}
	return false
}

var validPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if password, ok := fieldLevel.Field().Interface().(string); ok {
		return val.VerifyPassword(password)
	}
	return false
}

var validTimeOnly validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if timeOnly, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateTimeOnly(timeOnly)
	}
	return false
}

var validDateOnly validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if dateOnly, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateDateOnly(dateOnly)
	}
	return false
}

var validMoney validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if money, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateMoney(money)
	}
	return false
}

var validName validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if name, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateName(name)
	}
	return false
}

var validShortletSpace validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if space, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateShortletSpace(space)
	}
	return false
}

var validShortletType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if space, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateShortletType(space)
	}
	return false
}

var validEventLocationType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if eventLocationType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateLocationType(eventLocationType)
	}
	return false
}

var validEventDateType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if eventDateType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateEventDateType(eventDateType)
	}
	return false
}

var validEventTicketType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if eventTicketType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateEventTicketType(eventTicketType)
	}
	return false
}

var validEventTicketMainType validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if eventTicketMainType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateEventTicketMainType(eventTicketMainType)
	}
	return false
}

var validEventTicketLevel validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if eventTicketLevel, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateEventTicketLevel(eventTicketLevel)
	}
	return false
}

var validPropertySizeUnit validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if unit, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidatePropertyUnit(unit)
	}
	return false
}

var validDesTypes validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if desType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateDesTypes(desType)
	}
	return false
}

var validOptionExtraInfo validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if optionExtraInfoType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateOptionExtraInfo(optionExtraInfoType)
	}
	return false
}

var validUserProfileTypes validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if userProfileType, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateUserProfileType(userProfileType)
	}
	return false
}

var validReportOption validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if reportOption, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateReportOption(reportOption)
	}
	return false
}

var validGuestOption validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if guestOption, ok := fieldLevel.Field().Interface().([]string); ok {
		return val.ValidateGuestOptionTypes(guestOption)
	}
	return false
}

var validUserOptionCancel validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if reason, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateUserOptionCancelTypes(reason)
	}
	return false
}


var validUserEventCancel validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if reason, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateUserEventCancelTypes(reason)
	}
	return false
}

var validHostOptionCancel validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if reason, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateHostOptionCancelTypes(reason)
	}
	return false
}

var validHostEventCancel validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if reason, ok := fieldLevel.Field().Interface().(string); ok {
		return val.ValidateHostEventCancelTypes(reason)
	}
	return false
}