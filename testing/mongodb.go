package testing

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoOptions interface {
	GetOptions() *options.ClientOptions
}

var (
	client     *mongo.Client
	clientOnce sync.Once
)

const mongoTimeout = 5 * time.Second

func CreateMongoDBWithDBName(t testing.TB, optMaker MongoOptions) (*mongo.Client, string) {
	t.Helper()
	database := strings.ReplaceAll(t.Name(), "/", "_") + "_" + strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)

	clientOnce.Do(func() {
		var err error

		opt := optMaker.GetOptions()
		ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)

		defer cancel()

		client, err = mongo.Connect(ctx, opt)
		if err != nil {
			t.Fatalf("Failed to connect to Mongo. Error %v", err)
		}
	})

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
		defer cancel()

		if err := client.Database(database).Drop(ctx); err != nil {
			t.Errorf("Failed to drop Mongo database. Error: %v", err)
		}
	})

	return client, database
}
