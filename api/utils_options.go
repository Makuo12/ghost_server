package api

import (
	"encoding/json"

	db "github.com/makuo12/ghost_server/db/sqlc"
)

func StructToStringShortlet(shortlet db.Shortlet) (string, error) {
	bytes, err := json.Marshal(shortlet)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructShortlet(shortlet string) (db.Shortlet, error) {
	var shortletDB db.Shortlet
	err := json.Unmarshal([]byte(shortlet), &shortletDB)
	if err != nil {
		return db.Shortlet{}, err
	}
	return shortletDB, nil
}

func StructToStringLocation(location db.Location) (string, error) {
	bytes, err := json.Marshal(location)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructLocation(location string) (db.Location, error) {
	var locationDB db.Location
	err := json.Unmarshal([]byte(location), &locationDB)
	if err != nil {
		return db.Location{}, err
	}
	return locationDB, nil
}

func StructToStringSpaceArea(spaceArea []db.SpaceArea) (string, error) {
	bytes, err := json.Marshal(spaceArea)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructSpaceArea(spaceArea string) ([]db.SpaceArea, error) {
	var spaceAreaDB []db.SpaceArea
	err := json.Unmarshal([]byte(spaceArea), &spaceAreaDB)
	if err != nil {
		return []db.SpaceArea{}, err
	}
	return spaceAreaDB, nil
}

// Events

func StructToStringEventLocation(location []db.EventDateLocation) (string, error) {
	bytes, err := json.Marshal(location)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructEventLocation(location string) ([]db.EventDateLocation, error) {
	var locationDB []db.EventDateLocation
	err := json.Unmarshal([]byte(location), &locationDB)
	if err != nil {
		return []db.EventDateLocation{}, err
	}
	return locationDB, nil
}

func StructToStringEventDateDetail(eventDetail []db.EventDateDetail) (string, error) {
	bytes, err := json.Marshal(eventDetail)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructEventDateDetail(eventDetail string) ([]db.EventDateDetail, error) {
	var eventDetailDB []db.EventDateDetail
	err := json.Unmarshal([]byte(eventDetail), &eventDetailDB)
	if err != nil {
		return []db.EventDateDetail{}, err
	}
	return eventDetailDB, nil
}

func StructToStringEventDateTime(eventDateTime []db.EventDateTime) (string, error) {
	bytes, err := json.Marshal(eventDateTime)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructEventDateTime(eventDateTime string) ([]db.EventDateTime, error) {
	var eventDateTimeDB []db.EventDateTime
	err := json.Unmarshal([]byte(eventDateTime), &eventDateTimeDB)
	if err != nil {
		return []db.EventDateTime{}, err
	}
	return eventDateTimeDB, nil
}

func StructToStringEventInfo(eventInfo db.EventInfo) (string, error) {
	bytes, err := json.Marshal(eventInfo)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func StringToStructEventInfo(eventInfo string) (db.EventInfo, error) {
	var eventInfoDB db.EventInfo
	err := json.Unmarshal([]byte(eventInfo), &eventInfoDB)
	if err != nil {
		return db.EventInfo{}, err
	}
	return eventInfoDB, nil
}
