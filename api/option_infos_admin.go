package api

import (
	"context"
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"
)

func RemoveOptionInfoAdmin(ctx context.Context, server *Server) func() {
	return func() {
		names := []string{"ajibade", "grace", "richard", "ulc homes", "sylvia"}
		for _, name := range names {
			user, err := server.store.GetUserByFirstName(ctx, name)
			if err != nil {
				log.Printf("err at GetUserByFirstName: %v for name: %v\n", err, name)
				continue
			}
			options, err := server.store.GetHostOptionInfoByHost(ctx, user.ID)
			if err != nil {
				log.Printf("err at GetHostOptionInfoByHost: %v for name: %v\n", err, name)
				continue
			}
			for _, o := range options {
				charges, err := server.store.ListAllIDChargeOptionReference(ctx, o.OptionUserID)
				if err != nil {
					log.Printf("err at ListAllIDChargeOptionReference: %v for name: %v\n", err, name)
				}
				arg := db.DeleteOptionParams{
					OptionID:     o.ID,
					OptionUserID: o.OptionUserID,
					ChargeID:     charges,
				}
				_, err = server.store.DeleteOption(ctx, arg, server.Bucket)
				if err != nil {
					log.Printf("err at DeleteOption: %v for name: %v\n", err, name)
					return
				}
			}
		}
	}
}
