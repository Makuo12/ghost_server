package api

import (
	"context"
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"

	"github.com/google/uuid"
)

// SingleContextRoom setups up a single room using the context.Context
func SingleContextRoom(ctx context.Context, server *Server, userOne uuid.UUID, userTwo uuid.UUID, funcName string) (id uuid.UUID, err error) {
	// When we create a message we want to create a room is this user and the receiver doesn't have a room
	id, err = server.store.GetSingleRoomID(ctx, db.GetSingleRoomIDParams{
		UserOne:   userOne,
		UserTwo:   userTwo,
		UserOne_2: userTwo,
		UserTwo_2: userOne,
	})
	if err == nil {
		return
	} else {
		err = nil
		// Since there was an error we want to create the room
		id, err = server.store.CreateSingleRoom(ctx, db.CreateSingleRoomParams{
			UserOne: userOne,
			UserTwo: userTwo,
		})
		if err != nil {
			log.Printf("Room at SingleGinRoom store.CreateSingleRoom this means room was not created for %v and receiver %v funcName: %v, err: %v", userOne, userTwo, funcName, err.Error())
			return
		}
	}
	return
}
