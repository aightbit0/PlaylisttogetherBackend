package database

import (
	"database/sql"
	"fmt"
)

func CheckIfSongExits(db *sql.DB, uri string, plname string) bool {
	var login Playlist
	query := "select id from playlist WHERE uri = ? and playlistname = ?"
	row := db.QueryRow(query, uri, plname)
	err := row.Scan(&login.ID)

	switch err {
	case sql.ErrNoRows:
		return true
	case nil:
		return false
	default:
		fmt.Println(err)
	}

	return false
}
