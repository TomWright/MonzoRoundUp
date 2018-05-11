package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
	"fmt"
	"os"
)

func init() {
	goose.AddMigration(Up00003, Down00003)
}

func Up00003(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(""+
		"CREATE TABLE %s.userAccounts ("+
		"id   VARCHAR(100),"+
		"userId INT,"+
		"PRIMARY KEY (id),"+
		"INDEX (userId)"+
		");", os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}

func Down00003(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`DROP TABLE %s.userAccounts;`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}
