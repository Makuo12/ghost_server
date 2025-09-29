package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DBSource = "postgresql://postgres:Filler@localhost:5105/flex_server?sslmode=disable"
)

var testStore *SQLStore

func TestMain(m *testing.M) {
	connPool, err := pgxpool.New(context.Background(), DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
