package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
	"fmt"
	"os"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		CREATE DATABASE %s;
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		DROP DATABASE %s;
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}
