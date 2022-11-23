package testing

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gocql/gocql"
)


type ScyllaDBOptions interface {
	GetOptions() *gocql.ClusterConfig
}


func getScyllaDBConfig(t *testing.T, opt *gocql.ClusterConfig) *gocql.ClusterConfig {
	t.Helper()

	keyspace := "testing_keyspace"

	if value, exists := os.LookupEnv("SCYLLADB_HOSTS"); exists {
		opt.Hosts = strings.Split(value, ",")
	}

	if value, exists := os.LookupEnv("SCYLLADB_KEYSPACE"); exists {
		keyspace = value
	}


	if value, exists := os.LookupEnv("SCYLLADB_PORT"); exists {
		opt.Port, _ = strconv.Atoi(value)
	}

	opt.Keyspace = keyspace // FIXME: + "_" + strconv.FormatUint(uint64(rand.Int31n(100_000)), 10)

	return opt
}


func CreateScyllaDB(t *testing.T, optMaker ScyllaDBOptions) *gocql.Session {
	t.Helper()

	opt := optMaker.GetOptions()

	cluster := getScyllaDBConfig(t, opt)

	scyllaSession, err := cluster.CreateSession()
	if err != nil {
		t.Fatalf("Failed to create scylla session: %v", err)
	}

	t.Cleanup(func() {
		scyllaSession.Close()
	})

	return scyllaSession
}
