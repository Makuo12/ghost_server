package gapi

import (
	// "errors"
	"context"
	"fmt"
	"log"
	"path/filepath"

	//"log"
	// "net/http"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/pb"
	"github.com/makuo12/ghost_server/token"
	"github.com/makuo12/ghost_server/utils"

	// "strings"
	//"time"

	//
	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

var Context *gin.Context
var RedisClient *redis.Client

type Server struct {
	pb.UnimplementedFlexServiceServer
	tokenMaker token.Maker
	config     utils.Config
	store      *db.SQLStore
	Bucket     *storage.BucketHandle
}

var RedisContext = context.Background()

func NewServer(config utils.Config, store *db.SQLStore) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)

	}
	configFirebase := &firebase.Config{
		StorageBucket: config.FirebaseBucketName,
	}
	path, err := filepath.Abs(config.FirebaseAPIKeyFile)
	if err != nil {
		log.Printf("error you an error getting firebase json file %v\n", err.Error())
	}
	opt := option.WithCredentialsFile(path)
	appFirebase, err := firebase.NewApp(context.Background(), configFirebase, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := appFirebase.Storage(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		Bucket:     bucket,
	}
	// adr := []string{"redis_flex:6379"}
	// rdb := redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName: "Flex",
	// 	SentinelAddrs: adr,
	// 	DB: 0,
	// })
	rdb := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "", // no password set
		DB:       0,            // use default DB
	})
	RedisClient = rdb
	rdb.Ping(RedisContext)

	return server, nil
}
