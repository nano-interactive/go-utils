package testing

import (
	"database/sql"
	"fmt"
	"github.com/tanimutomo/sqlfile"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TimescaledbptionsData struct {
	Host          string
	Username      string
	Password      string
	DB            string
	Schema        string
	SSLMode       string
	Port          int
	MigrationFile string
}

type TimescaledbOptions interface {
	GetOptions() TimescaledbptionsData
}

var (
	tsdbClient      *sql.DB
	timescaledbOnce sync.Once
)

func CreateTimescaledb(t testing.TB, optMaker TimescaledbOptions) (*sql.DB, string) {
	opts := optMaker.GetOptions()
	dbName := strings.ReplaceAll(t.Name(), "/", "_") + "_" + strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	dbName = strings.ToLower(dbName)

	var err error

	timescaledbOnce.Do(func() {
		tsdbClient, err = sql.Open("pgx", fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Port,
			opts.DB,
		))
		if err != nil {
			t.Error("Failed to connect to timescaledb", err)
			t.FailNow()
		}
	})

	if _, err = tsdbClient.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
		t.Errorf("Failed to create database with name %s: %v", dbName, err)
		t.FailNow()
	}
	var client *sql.DB

	t.Cleanup(func() {

		if client != nil {
			_ = client.Close()
		}
		_, err = tsdbClient.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
		if err != nil {
			t.Errorf("failed to delete database %s: %v", dbName, err)
		}
	})

	if _, err = tsdbClient.Exec(fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s;", dbName, opts.Username)); err != nil {
		t.Errorf("Failed to grant privileges for database %s to user %s: %v", dbName, opts.Username, err)
		t.FailNow()
	}

	// switch databases
	client, err = sql.Open("pgx", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		dbName,
	))
	if err != nil {
		t.Error("Failed to connect to timescaledb", err)
		t.FailNow()
	}

	s := sqlfile.New()
	err = s.File(opts.MigrationFile)
	if err != nil {
		t.Errorf("failed to read Timescaledb migrations file (DIR=%s): %v", opts.MigrationFile, err)
		t.FailNow()
	}

	if _, err = s.Exec(client); err != nil {
		t.Errorf("failed to exec Timescaledb migrations file (DIR=%s): %v", opts.MigrationFile, err)
		t.FailNow()
	}

	return client, dbName

}
