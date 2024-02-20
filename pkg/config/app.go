package config

import (
	"context"
	"fmt"
	"log"

	"github.com/abhik-99/passwordless-login/pkg/utils"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Db       *mongo.Database
	MongoCtx context.Context

	client *mongo.Client

	Rdb      *redis.Client
	RedisCtx context.Context
)

func init() {
	MongoCtx = context.Background()
	client, err := mongo.Connect(MongoCtx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:27017/?retryWrites=true&w=majority", utils.GetENV("DBUSER"), utils.GetENV("DBPASS"))))
	if err != nil {
		// panic(err)
		log.Println(err)
	}

	if err = client.Ping(MongoCtx, readpref.Primary()); err == nil {
		log.Print("Connection to DB Successful")
	} else {
		log.Println("ERROR while pinging DB")
		log.Panic(err)
		return
	}

	Db = client.Database("passwordless-auth")
	Db.CreateCollection(MongoCtx, "user-collection")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetENV("REDISADDR"),
		Password: utils.GetENV("REDISPASS"),
		DB:       0, // use default DB
	})

	RedisCtx = context.Background()
}

func Disconnect() {
	client.Disconnect(MongoCtx)

	Rdb.Close()
}
