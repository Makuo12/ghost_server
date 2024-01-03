package main

import (
	"context"
	"flex_server/api"
	db "flex_server/db/sqlc"

	//"flex_server/gapi"
	//"flex_server/pb"
	"flex_server/utils"
	"log"

	//"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "time/tzdata"
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
