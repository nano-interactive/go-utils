package testing

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tanimutomo/sqlfile"
)

type MySQLOptionsData struct {
	RootUser      string
	RootPass      string
	Host          string
	Net           string
	DB            string
	User          string
	Password      string
	MigrationFile string
}

type MySQLOptions interface {
	GetOptions() MySQLOptionsData
}

var (
	mysqlClient     *sql.DB
	mysqlClientOnce sync.Once
)

func CreateMySQL(t testing.TB, optMaker MySQLOptions) (*sql.DB, string) {
	opts := optMaker.GetOptions()

	mysqlClientOnce.Do(func() {
		var err error
		mysqlClient, err = sql.Open(
			"mysql",
			fmt.Sprintf(
				"%s:%s@%s(%s)/",
				opts.RootUser,
				opts.RootPass,
				opts.Net,
				opts.Host,
			),
		)
		if err != nil {
			t.Errorf("Failed to open new SQL connection: %v", err)
			t.FailNow()
		}
	})

	dbName := strings.ReplaceAll(t.Name(), "/", "_") + "_" + strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)

	_, err := mysqlClient.Exec(
		fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName),
	)
	if err != nil {
		t.Errorf("Failed to create database with name %s: %v", dbName, err)
		t.FailNow()
	}

	var client *sql.DB

	t.Cleanup(func() {
		if client != nil {
			_ = client.Close()
		}
		_, err = mysqlClient.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
		if err != nil {
			t.Errorf("failed to delete database %s: %v", dbName, err)
		}
	})

	_, err = mysqlClient.Exec(
		fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' WITH GRANT OPTION;", dbName, opts.User),
	)
	if err != nil {
		t.Errorf("Failed to grant privileges for database %s to user %s: %v", dbName, opts.User, err)
		t.FailNow()
	}

	client, err = sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@%s(%s)/%s",
			opts.User,
			opts.Password,
			opts.Net,
			opts.Host,
			dbName,
		),
	)
	if err != nil {
		t.Errorf("failed to open connection to database: %s %v", dbName, err)
		t.FailNow()
	}

	s := sqlfile.New()
	err = s.File(opts.MigrationFile)
	if err != nil {
		t.Errorf("failed to read MySQL migrations file (DIR=%s): %v", opts.MigrationFile, err)
		t.FailNow()
	}

	if _, err = s.Exec(client); err != nil {
		t.Errorf("failed to read MySQL migrations file (DIR=%s): %v", opts.MigrationFile, err)
		t.FailNow()
	}

	return client, dbName
}
