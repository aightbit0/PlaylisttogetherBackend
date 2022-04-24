package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"playlisttogether/backend/utils"
	"strings"
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

func CheckIfSongDislike(db *sql.DB, id int, user string) (int, string, error) {

	var login Playlist
	query := "select dislike, disliker from playlist WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&login.Dislike, &login.Disliker)

	dislikes := login.Dislike

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return 0, "", errors.New("no rows returned")
	case nil:
		index := -1

		s := strings.Split(login.Disliker, ",")
		for i := 0; i < len(s); i++ {
			if strings.ToLower(s[i]) == strings.ToLower(user) {
				index = i
			}
		}

		if index >= 0 {
			s = utils.RemoveIndex(s, index)
			dislikes = dislikes - 1
		} else {
			s = append(s, user)
			dislikes = dislikes + 1
		}

		fmt.Println(len(s))
		fmt.Println(s)
		var r []string
		for _, str := range s {
			if str != "" {
				r = append(r, str)
			}
		}

		return dislikes, strings.Join(r, ","), nil

	default:
		fmt.Println(err)
	}

	return 0, "", err
}

func SelectBucket(db *sql.DB, user string, plname string) ([]byte, error) {
	var array []Playlist
	query := "select * from playlist WHERE user = ? and playlistname = ?"
	rows, err := db.Query(query, user, plname)
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Playlist
		err := rows.Scan(&tag.ID, &tag.User, &tag.Songname, &tag.Artist, &tag.Url, &tag.Picture, &tag.Dislike, &tag.Uri, &tag.Playlist, &tag.Disliker, &tag.PlaylistName)
		if err != nil {
			fmt.Println("Failed Scan")
			return nil, err
		}

		array = append(array, tag)
	}
	err = rows.Err()

	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}

	amount, err := SelectAmount(db, user, plname)

	if err != nil {
		fmt.Println("failed to get Amount")
		return nil, err
	}

	var returnVal Bucket
	returnVal.Data = array
	returnVal.Amount = amount

	end, err := json.Marshal(returnVal)
	if err != nil {
		fmt.Println(("Failes Marshal"))
		return nil, err
	}

	return end, nil

}

func SelectPlaylist(db *sql.DB, plname string) ([]byte, error) {
	var array []Playlist
	query := "select * from playlist where playlist = true and playlistname = ?"
	rows, err := db.Query(query, plname)
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Playlist
		err := rows.Scan(&tag.ID, &tag.User, &tag.Songname, &tag.Artist, &tag.Url, &tag.Picture, &tag.Dislike, &tag.Uri, &tag.Playlist, &tag.Disliker, &tag.PlaylistName)
		if err != nil {
			fmt.Println("Failed Scan")
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
		fmt.Println(("Failes Marshal"))
		return nil, err
	}

	return end, nil

}

func SelectUris(db *sql.DB, user string, plname string) ([]byte, error) {
	var uris []string
	query := "select uri from playlist where playlist = true and playlistname = ?"
	rows, err := db.Query(query, plname)
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

	end, err := json.Marshal(uris)
	if err != nil {
		return nil, err
	}

	return end, nil

}

func DislikeSong(db *sql.DB, id int, user string, plname string) (bool, error) {

	amount, err := SelectAmountOfUsers(db, plname)

	if err != nil {
		return false, err
	}

	dislikes, disliker, err := CheckIfSongDislike(db, id, user)

	fmt.Println(disliker)

	if err != nil {
		fmt.Println("Error")
		return false, err
	}

	stmt, e := db.Prepare("update playlist set dislike =? , disliker = ? where id=?")
	if e != nil {
		log.Printf("%v", e.Error())
	}
	// execute
	res, e := stmt.Exec(dislikes, disliker, id)
	if e != nil {
		log.Printf("%v", e.Error())
	}

	a2, e := res.RowsAffected()
	if e != nil {
		log.Printf("%v", e.Error())
	}

	if a2 == 0 {
		fmt.Println("DislikeSong 0 rows affected")
		return false, errors.New(" 0 rows returned")
	}

	faktor := amount / 2

	if dislikes >= faktor {
		//fmt.Println("DELETE SONG:")

		number, err := DeleteSongSystem(db, id)
		if err != nil || number == 0 {
			return false, nil
		}

		return true, nil
	}

	return false, nil

}

// only system can Delete a Song without user.
func DeleteSongSystem(db *sql.DB, id int) (int64, error) {
	stmt, e := db.Prepare("delete from playlist where id=?")
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}
	// execute
	res, e := stmt.Exec(id)
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}

	a, e := res.RowsAffected()
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}

	return a, nil
}

// DeleteSong for User
func DeleteSong(db *sql.DB, id int, user string) (int64, error) {
	stmt, e := db.Prepare("delete from playlist where id=? and user=?")
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}
	// execute
	res, e := stmt.Exec(id, user)
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}

	a, e := res.RowsAffected()
	if e != nil {
		log.Printf("%v", e.Error())
		return 0, e
	}

	return a, nil
}

func Insert(db *sql.DB, obj Playlist) (int64, error) {

	stmt, err := db.Prepare("INSERT INTO playlist SET user = ?, songname = ?, artist = ?, url = ?, picture = ? ,dislike = ?, uri = ?, playlist = ?, disliker = '', playlistname = ?")
	if err != nil {
		log.Printf("%v", err.Error())
		return 0, err
	}
	res, err := stmt.Exec(obj.User, obj.Songname, obj.Artist, obj.Url, obj.Picture, 0, obj.Uri, true, obj.PlaylistName)
	if err != nil {
		log.Printf("%v", err.Error())
		return 0, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		log.Printf("%v", err.Error())
		return 0, err
	}

	return lid, nil
}
