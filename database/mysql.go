package database

import (
	"database/sql"
	"fmt"
	"playlisttogether/backend/config"

	_ "github.com/go-sql-driver/mysql"
)

type db struct {
	*sql.DB
}

func NewDB(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open(config.Type, config.User+":"+config.Password+"@tcp("+config.URL+":"+config.Port+")/"+config.Dbname)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("no connection")
		return nil, err
	}

	fmt.Println("connection stable")
	return db, nil
}
