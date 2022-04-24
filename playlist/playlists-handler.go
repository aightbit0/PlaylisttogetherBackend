package playlist

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"playlisttogether/backend/database"
)

func (p *PlaylistCredentials) GetPlaylists(db *sql.DB) http.HandlerFunc {
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

		end, err := database.SelectPlaylists(db, t.User)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Select has failed"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) GetSongUris(db *sql.DB) http.HandlerFunc {
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
		end, err := database.SelectUris(db, t.User, t.PlaylistName)
		if err != nil {
			fmt.Println("failed to Select")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("failed to Select"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) UpdatePlaylist(db *sql.DB) http.HandlerFunc {
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

		database.UpdatePlaylist(db, t.User, t.PlaylistName, t.PlaylistUrl)

		end, err := json.Marshal("sucess")
		if err != nil {
			fmt.Println(("FAILED"))
		}

		rw.Write(end)
	}
}

func CreatePlaylist(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var t database.Playlists
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		res, err := database.CheckIfPlaylistExits(db, t.PlaylistName)

		if err != nil {
			fmt.Println("failed select playlist")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("failed select playlist"))
		}

		if res {
			end, err := json.Marshal("exists")
			if err != nil {
				fmt.Println("failed to marshal")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("JSON not valid"))
			}
			rw.Write(end)
		}

		playlist, err := database.CreatePlaylist(db, t)

		if err != nil || !playlist {
			fmt.Println("failed create playlist")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("failed create playlist"))
		}

		end, err := json.Marshal("success")
		if err != nil {
			fmt.Println("failed to marshal")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		rw.Write(end)

	}
}
