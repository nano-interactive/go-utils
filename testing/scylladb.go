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

func getScyllaDBConfig(t testing.TB, opt *gocql.ClusterConfig) *gocql.ClusterConfig {
	t.Helper()

	if value, exists := os.LookupEnv("SCYLLADB_HOSTS"); exists {
		opt.Hosts = strings.Split(value, ",")
	}

	if value, exists := os.LookupEnv("SCYLLADB_KEYSPACE"); exists {
		opt.Keyspace = value
	}

	if value, exists := os.LookupEnv("SCYLLADB_PORT"); exists {
		opt.Port, _ = strconv.Atoi(value)
	}

	return opt
}

func CreateScyllaDB(t testing.TB, optMaker ScyllaDBOptions) *gocql.Session {
	t.Helper()

	opt := optMaker.GetOptions()

	cluster := getScyllaDBConfig(t, opt)

	scyllaSession, err := cluster.CreateSession()
	if err != nil {
		t.Errorf("Failed to create scylla session: %v", err)
		t.FailNow()
	}

	t.Cleanup(func() {
		scyllaSession.Close()
	})

	return scyllaSession
}
