package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"playlisttogether/backend/config"
	"playlisttogether/backend/utils"
	"strconv"
	"strings"

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

func CheckIfUserAuth(db *sql.DB, uname string, pword string) int {
	var login PlaylistUser
	query := "select id, user, password, active from playlistuser WHERE user = ?"
	row := db.QueryRow(query, uname)
	err := row.Scan(&login.ID, &login.Name, &login.Password, &login.Active)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return -1
	case nil:
		if utils.GetMD5Hash(pword) == login.Password {
			return login.ID
		}
		return 0

	default:
		fmt.Println(err)
	}

	return -1
}

//UPDATE
func SetOnlineStatus(db *sql.DB, user string, status bool) {
	stmt, e := db.Prepare("update playlistuser set active=? where user=?")
	if e != nil {
		log.Printf("%v", e.Error())
	}
	// execute
	res, e := stmt.Exec(status, user)
	if e != nil {
		log.Printf("%v", e.Error())
	}

	a, e := res.RowsAffected()
	if e != nil {
		log.Printf("%v", e.Error())
	}

	if a == 0 {
		fmt.Println("setOnlineStatus 0 rows affected")
	}
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

func SelectUsers(db *sql.DB, user string) ([]byte, error) {
	var users PlaylistUser
	var allUsers []PlaylistUser
	query := "select id, user from playlistuser WHERE user NOT LIKE ?"
	rows, err := db.Query(query, user)
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&users.ID, &users.Name)
		if err != nil {
			fmt.Println("Failed Select")
			return nil, err
		}
		userFromTable := PlaylistUser{Name: users.Name, ID: users.ID}
		allUsers = append(allUsers, userFromTable)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Failed Select")
		return nil, err
	}

	end, err := json.Marshal(allUsers)
	if err != nil {
		fmt.Println("Failed marahal")
		return nil, err
	}

	return end, nil

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

func DislikeSong(db *sql.DB, id int, user string, plname string) (bool, error) {

	amount, err := SelectAmountOfUsers(db, plname)

	if err != nil {
		return false, err
	}

	dislikes, disliker, err := CheckIfSongDislike(db, id, user)

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
		fmt.Println("DELETE SONG:")

		number, err := DeleteSongSystem(db, id)
		if err != nil || number == 0 {
			return false, nil
		}

		return true, nil
	}

	return false, nil

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

		return dislikes, strings.Join(s, ","), nil

	default:
		fmt.Println(err)
	}

	return 0, "", err
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
