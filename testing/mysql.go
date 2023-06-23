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
		t.Fatalf("Failed to open new SQL connection: %v", err)
	}

	dbName := opts.DB + "_" + strings.ReplaceAll(utils.RandomString(10), "-", "_")

	_, err = client.Exec(
		fmt.Sprintf("CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName),
	)

	if err != nil {
		t.Fatal(err)
	}

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
			t.Logf("failed to drop database: %s %v", dbName, err)
		}

		defer client.Close()

		_, err = client.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
		if err != nil {
			t.Logf("failed to delete database %s: %v", dbName, err)
		}
	})

	t.Logf("Username=%s Password=%s Host=%s DB=%s", opts.User, opts.Password, opts.Host, dbName)

	_, err = client.Exec(
		fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' WITH GRANT OPTION;", dbName, opts.User),
	)

	if err != nil {
		t.Fatal(err)
	}

	if err = client.Close(); err != nil {
		t.Logf("failed to close connection: %v", err)
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
		t.Fatalf("failed to open connection to database: %s %v", dbName, err)
	}

	s := sqlfile.New()
	err = s.File(opts.MigrationFile)
	if err != nil {
		t.Fatalf("failed to read MySQL migrations file (DIR=%s): %v", opts.MigrationFile, err)
	}

	if _, err = s.Exec(client); err != nil {
		t.Fatalf("failed to read MySQL migrations file (DIR=%s): %v", opts.MigrationFile, err)
	}

	t.Cleanup(func() {
		if err = client.Close(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	})

	return client, dbName
}
