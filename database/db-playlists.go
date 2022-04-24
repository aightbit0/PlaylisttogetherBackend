package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

func SelectAmount(db *sql.DB, user string, plname string) (int, error) {
	var amount int
	query := "select amount from playlists where user = ? and playlistname = ?"
	row := db.QueryRow(query, user, plname)
	err := row.Scan(&amount)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0, nil
	case nil:
		return amount, nil

	default:
		fmt.Println(err)
		return 0, err
	}
}

func SelectPlaylists(db *sql.DB, user string) ([]byte, error) {
	var array []Playlists
	query := "select * from playlists where user = ?"
	rows, err := db.Query(query, user)
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag Playlists
		err := rows.Scan(&tag.ID, &tag.User, &tag.PlaylistName, &tag.PlaylistUrl, &tag.PlaylistID, &tag.Amount, &tag.Creator)
		if err != nil {
			fmt.Println("Failed Select")
			return nil, err
		}

		array = append(array, tag)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}

	end, err := json.Marshal(array)
	if err != nil {
		fmt.Println("Failed marshal")
		return nil, err
	}

	return end, nil

}

//SelectAmountOfUsers
func SelectAmountOfUsers(db *sql.DB, plname string) (int, error) {
	var amount string
	query := "SELECT COUNT(user) as anzahl FROM playlists where playlistname = ?"
	row := db.QueryRow(query, plname)
	err := row.Scan(&amount)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0, errors.New("no rows returned")
	case nil:

		am, err := strconv.Atoi(amount)
		if err != nil {
			return 0, err
		}
		return am, nil

	default:
		fmt.Println(err)
	}

	return 0, err
}

func CheckIfPlaylistExits(db *sql.DB, plname string) (bool, error) {

	var login Playlists
	query := "select id from playlists WHERE playlistname = ?"
	row := db.QueryRow(query, plname)
	err := row.Scan(&login.ID)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return false, nil
	case nil:

		return true, nil

	default:
		fmt.Println(err)
	}

	return false, err
}

func CreatePlaylist(db *sql.DB, obj Playlists) (bool, error) {

	for i := 0; i < len(obj.Users); i++ {

		stmt, err := db.Prepare("INSERT INTO playlists SET user = ?, playlistname = ?, playlisturl = ?, playlistid = ?, amount = ? , creator = ?")
		if err != nil {
			log.Printf("%v", err.Error())
			return false, err
		}
		res, err := stmt.Exec(obj.Users[i].Label, obj.PlaylistName, "", "", obj.Amount, obj.User)
		if err != nil {
			log.Printf("%v", err.Error())
			return false, err
		}

		lid, err := res.LastInsertId()
		if err != nil {
			log.Printf("%v", err.Error())
			return false, err
		}

		if lid == 0 {
			fmt.Println("Insert 0 rows affected")
			return false, errors.New("0 rows affected")
		}

	}

	return true, nil

}

func UpdatePlaylist(db *sql.DB, user string, plname string, plurl string) error {

	stmt, e := db.Prepare("update playlists set playlisturl =? where creator=? and playlistname = ?")
	if e != nil {
		log.Printf("%v", e.Error())
		return e
	}
	// execute
	res, e := stmt.Exec(plurl, user, plname)
	if e != nil {
		log.Printf("%v", e.Error())
		return e
	}

	a, e := res.RowsAffected()
	if e != nil {
		log.Printf("%v", e.Error())
		return e
	}
	if a == 0 {
		fmt.Println("updatePlaylist 0 rows affected")
		return errors.New(" 0 rows affected")
	}
	return nil
}
