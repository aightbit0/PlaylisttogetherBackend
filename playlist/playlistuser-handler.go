package playlist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"playlisttogether/backend/database"
)

func (p *PlaylistCredentials) GetUsers(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if *p.headerFlag {
			(rw).Header().Set("Access-Control-Allow-Origin", "*")
			(rw).Header().Set("Access-Control-Allow-Methods", "POST")
			(rw).Header().Set("Access-Control-Allow-Headers", "*")
		}

		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var t Auth
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		end, err := database.SelectUsers(db, t.User)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Select has failed"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) Logout(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if *p.headerFlag {
			(rw).Header().Set("Access-Control-Allow-Origin", "*")
			(rw).Header().Set("Access-Control-Allow-Methods", "POST")
			(rw).Header().Set("Access-Control-Allow-Headers", "*")
		}
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var t Auth
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("Failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		if t.User != "" {
			database.SetOnlineStatus(db, t.User, false)
		}

		rw.Write([]byte("logged out"))
	}
}
