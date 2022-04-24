package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"playlisttogether/backend/utils"
)

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
