package testing

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/nano-interactive/go-utils"
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

func CreateMySQL(t testing.TB, optMaker MySQLOptions) (*sql.DB, string) {
	opts := optMaker.GetOptions()

	client, err := sql.Open(
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

	dbName := opts.DB + "_" + strings.ReplaceAll(utils.RandomString(10), "-", "_")

	_, err = client.Exec(
		fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName),
	)

	if err != nil {
		t.Errorf("Failed to create database with name %s: %v", dbName, err)
		t.FailNow()
	}

	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close connection: %v", err)
		}
	})

	t.Cleanup(func() {
		client, err = sql.Open(
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
			t.Errorf("failed to drop database: %s %v", dbName, err)
			return
		}

		defer client.Close()

		_, err = client.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
		if err != nil {
			t.Errorf("failed to delete database %s: %v", dbName, err)
		}
	})

	_, err = client.Exec(
		fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' WITH GRANT OPTION;", dbName, opts.User),
	)

	if err != nil {
		t.Errorf("Failed to grant privileges for database %s to user %s: %v", dbName, opts.User, err)
		t.FailNow()
	}

	if err = client.Close(); err != nil {
		t.Errorf("failed to close connection: %v", err)
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
