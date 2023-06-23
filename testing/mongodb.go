package testing

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions interface {
	GetOptions() *options.ClientOptions
}

func getMongoDBConfig(t testing.TB, opt *options.ClientOptions) (*options.ClientOptions, string) {
	t.Helper()
	return opt, "testing_database_" + strconv.FormatUint(uint64(rand.Int31n(100_000)), 10)
}

func CreateMongoDBWithDBName(t testing.TB, optMaker MongoOptions) (*mongo.Client, string) {
	t.Helper()

	opt := optMaker.GetOptions()
	cfg, database := getMongoDBConfig(t, opt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to Mongo. Error %v", err)
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Database(database).Drop(ctx); err != nil {
			t.Fatalf("Failed to drop Mongo database. Error: %v", err)
		}

		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("Failed to drop Mongo connection. Error %v", err)
		}
	})

	return client, database
}
