package database

import (
	"database/sql"
)

func SelectPictures(db *sql.DB, folder string, from int, to int) ([]string, error) {
	var uris []string
	query := "select path from pictures where folder = ? Limit ?,?"
	rows, err := db.Query(query, folder, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Playlist
		err := rows.Scan(&tag.Uri)
		if err != nil {
			return nil, err
		}
		uris = append(uris, tag.Uri)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return uris, nil

}
