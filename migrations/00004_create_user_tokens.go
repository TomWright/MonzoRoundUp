package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
	"os"
	"fmt"
)

func init() {
	goose.AddMigration(Up00004, Down00004)
}

func Up00004(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		CREATE TABLE %s.userTokens (
		id INT AUTO_INCREMENT,
		accessToken   VARCHAR(255),
		clientId   VARCHAR(255),
		refreshToken   VARCHAR(255),
		tokenType   ENUM('BEARER'),
		monzoUserId   VARCHAR(255),
		userId   INT,
		expiresAt   DATETIME,
		PRIMARY KEY (id),
		INDEX(userId)
		);
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}

func Down00004(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`
		DROP TABLE %s.userTokens;
	`, os.Getenv("MYSQL_DB")))
	if err != nil {
		return err
	}
	return nil
}
