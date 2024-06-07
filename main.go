package main

import (
	"context"

	"github.com/makuo12/ghost_server/api"
	db "github.com/makuo12/ghost_server/db/sqlc"

	//"github.com/makuo12/ghost_server/gapi"
	//"github.com/makuo12/ghost_server/pb"
	"log"

	"github.com/makuo12/ghost_server/utils"

	//"net"

	_ "time/tzdata"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
)

func main() {
	go api.H.Run()
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not access env variables", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err.Error())
	}
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
		return
	}

	//api.H.Server = server
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err.Error())
	}
	//runGrpcServer(config, store)
	//runGinServer(config, store)
}

//func runGrpcServer(config utils.Config, store *db.SQLStore){
//	server, err := gapi.NewServer(config, store)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}

//	grpcServer := grpc.NewServer()
//	pb.RegisterFlexServiceServer(grpcServer, server)
//	reflection.Register(grpcServer)
//	listener, err := net.Listen("tcp", config.GRPCServerAddress)
//	if err != nil {
//		log.Fatal("Cannot start server", err.Error())
//	}
//	log.Printf("start grpc server at %s", listener.Addr().String())
//	err = grpcServer.Serve(listener)
//	if err != nil {
//		log.Fatal("Cannot start grpc server")
//	}
//}
//func runGinServer(config utils.Config, store *db.SQLStore) {
//	server, err := api.NewServer(config, store)
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	// api.H.Server = server
//	err = server.Start(config.HTTPServerAddress)
//	if err != nil {
//		log.Fatal("Cannot start server", err.Error())
//	}
//}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}

