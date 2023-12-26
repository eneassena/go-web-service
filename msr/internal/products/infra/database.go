package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func settingsDatabaseCennection() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("USER_DB"),
		os.Getenv("PASSWD_DB"),
		os.Getenv("HOST_DB"),
		os.Getenv("PORT_DB"),
		os.Getenv("NAME_DB"),
	)
}

func Connect() *sql.DB {
	conn, err := sql.Open("mysql", settingsDatabaseCennection())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
