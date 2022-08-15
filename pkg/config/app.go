package config

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Db  *mongo.Database
	Ctx context.Context

	client *mongo.Client
	cancel context.CancelFunc
)

func Connect() {
	var myEnv map[string]string
	myEnv, _ = godotenv.Read("../../.env")
	Ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.w7ovegb.mongodb.net/?retryWrites=true&w=majority", myEnv["DBUser"], myEnv["DBPass"])))
	if err != nil {
		// panic(err)
		fmt.Println(err)
	}
	Db = client.Database("passwordless-auth")
}

func Disconnect() {
	client.Disconnect(Ctx)
	cancel()
}
