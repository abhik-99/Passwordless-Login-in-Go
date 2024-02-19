package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/abhik-99/passwordless-login/pkg/utils"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Db       *mongo.Database
	MongoCtx context.Context

	client      *mongo.Client
	mongoCancel context.CancelFunc

	Rdb         *redis.Client
	RedisCtx    context.Context
	redisCancel context.CancelFunc
)

func init() {
	MongoCtx, mongoCancel = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(MongoCtx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@localhost:27017/?retryWrites=true&w=majority", utils.GetENV("DBUSER"), utils.GetENV("DBPASS"))).SetConnectTimeout(30*time.Second))
	if err != nil {
		// panic(err)
		log.Println(err)
	}
	Db = client.Database("passwordless-auth")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     utils.GetENV("REDISADDR"),
		Password: utils.GetENV("REDISPASS"),
		DB:       0, // use default DB
	})

	RedisCtx, redisCancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func Disconnect() {
	client.Disconnect(MongoCtx)
	mongoCancel()

	Rdb.Close()
	redisCancel()
}
