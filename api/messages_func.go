package api

import (
	"fmt"
	"log"
	"sort"
	"time"
)

// Functions
//func GenerateUniqueMessageContacts(redisData, resData []MessageContactItem) []MessageContactItem {
//	// Concatenate the two slices into a single slice.
//	allData := append(redisData, resData...)

//	// Create a map to keep track of unique ConnectedUserID.
//	uniqueUsers := make(map[string]struct{})

//	// Filter out duplicate ConnectedUserID and update the uniqueUsers map.
//	var uniqueMessages []MessageContactItem
//	for _, message := range allData {
//		if _, ok := uniqueUsers[message.ConnectedUserID]; !ok {
//			uniqueUsers[message.ConnectedUserID] = struct{}{}
//			uniqueMessages = append(uniqueMessages, message)
//		}
//	}

//	// Sort the uniqueMessages slice in descending order based on LastMessageTime.
//	sort.SliceStable(uniqueMessages, func(i, j int) bool {
//		return uniqueMessages[i].LastMessageTime > uniqueMessages[j].LastMessageTime
//	})

//	return uniqueMessages
//}

func GenerateUniqueMessageContacts(redisData, resData []MessageContactItem) []MessageContactItem {
	// Concatenate the two slices into a single slice.
	allData := append(redisData, resData...)
	// Sort the uniqueMessages slice in descending order based on LastMessageTime.
	sort.SliceStable(allData, func(i, j int) bool {
		return allData[i].LastMessageTime < allData[j].LastMessageTime
	})
	// Create a map to keep track of unique ConnectedUserID.
	users := make(map[string]MessageContactItem)

	for index, msg := range allData {

		newMsg := MessageContactItem{
			MsgID:                      msg.MsgID,
			ConnectedUserID:            msg.ConnectedUserID,
			FirstName:                  msg.FirstName,
			MainImage:                  msg.MainImage,
			LastMessage:                msg.LastMessage,
			LastMessageTime:            msg.LastMessageTime,
			UnreadMessageCount:         users[msg.ConnectedUserID].UnreadMessageCount + msg.UnreadMessageCount,
			UnreadUserRequestCount:     users[msg.ConnectedUserID].UnreadUserRequestCount + msg.UnreadUserRequestCount,
			UnreadUserCancelCount:      users[msg.ConnectedUserID].UnreadUserCancelCount + msg.UnreadUserCancelCount,
			UnreadHostCancelCount:      users[msg.ConnectedUserID].UnreadHostCancelCount + msg.UnreadHostCancelCount,
			UnreadHostChangeDatesCount: users[msg.ConnectedUserID].UnreadHostChangeDatesCount + msg.UnreadHostChangeDatesCount,
		}
		log.Printf("%v newMsg %v\n", index, newMsg)
		users[msg.ConnectedUserID] = newMsg
	}

	var data []MessageContactItem

	for _, u := range users {
		data = append(data, u)
	}

	return data
}

func GenerateUniqueMessages(redisData, resData []MessageItem) []MessageItem {
	// Concatenate the two slices into a single slice.
	allData := append(redisData, resData...)

	// Create a map to keep track of unique IDs.
	uniqueIDs := make(map[string]struct{})

	// Filter out duplicate IDs and update the uniqueIDs map.
	var uniqueMessages []MessageItem
	for _, message := range allData {
		if _, ok := uniqueIDs[message.ID]; !ok {
			uniqueIDs[message.ID] = struct{}{}
			uniqueMessages = append(uniqueMessages, message)
		}
	}

	// Sort the uniqueMessages slice in descending order based on CreatedAt.
	sort.SliceStable(uniqueMessages, func(i, j int) bool {
		return uniqueMessages[i].CreatedAt > uniqueMessages[j].CreatedAt
	})

	return uniqueMessages
}

func FindLatestLastMessageTime(items []MessageContactItem) (string, error) {
	if len(items) == 0 {
		return "", fmt.Errorf("empty list")
	}

	// Assume the first item has the latest LastMessageTime
	latestTime, err := time.Parse("2006-01-02T15:04:05", items[0].LastMessageTime)
	if err != nil {
		return "", err
	}

	// Iterate through the rest of the items and find the latest LastMessageTime
	for _, item := range items[1:] {
		timeObj, err := time.Parse("2006-01-02T15:04:05", item.LastMessageTime)
		if err != nil {
			return "", err
		}

		if timeObj.After(latestTime) {
			latestTime = timeObj
		}
	}

	// Format the latest LastMessageTime as a string
	latestTimeString := latestTime.Format("2006-01-02T15:04:05")
	return latestTimeString, nil
}
