package config

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

func Connect() {
	var myEnv map[string]string
	myEnv, _ = godotenv.Read("../../.env")
	MongoCtx, mongoCancel = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(MongoCtx, options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.w7ovegb.mongodb.net/?retryWrites=true&w=majority", myEnv["DBUser"], myEnv["DBPass"])))
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
	Db = client.Database("passwordless-auth")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: myEnv["RedisPass"],
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
