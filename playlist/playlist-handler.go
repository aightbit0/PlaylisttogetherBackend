package playlist

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"playlisttogether/backend/database"
	"strings"

	"github.com/gorilla/mux"
)

type PlaylistCredentials struct {
	hmacSampleSecret []byte
	headerFlag       *bool
}

func NewPlaylistInstance(secret []byte, flag *bool, expireTime int) PlaylistCredentials {
	playlistInstance := PlaylistCredentials{
		hmacSampleSecret: secret,
		headerFlag:       flag,
	}

	fmt.Println("new Playlist Instance created")

	return playlistInstance
}

// Middleware checks if User is Auth
func (p *PlaylistCredentials) JwtMiddleware(db *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if *p.headerFlag {
				(w).Header().Set("Access-Control-Allow-Origin", "*")
				(w).Header().Set("Access-Control-Allow-Methods", "*")
				(w).Header().Set("Access-Control-Allow-Headers", "*")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

			if len(authHeader) != 2 {
				fmt.Println("Malformed token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))

			} else {

				claims, err := p.checkJWT(authHeader[1])
				if err != nil {
					fmt.Println("unguilty Token")
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}
				ctx := context.WithValue(r.Context(), "user", claims)
				uid := database.CheckIfUserAuth(db, claims["name"].(string), "")
				if uid == 0 {
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}
			}
		})
	}
}

func (p *PlaylistCredentials) Login(db *sql.DB) http.HandlerFunc {
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
		var t database.PlaylistUser
		err := decoder.Decode(&t)

		if err != nil {
			fmt.Println("Failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		isAuth := database.CheckIfUserAuth(db, t.Name, t.Password)

		if isAuth > 0 {
			resonse := p.createJWT(t.Name)
			database.SetOnlineStatus(db, t.Name, true)

			end, err := json.Marshal(resonse)
			if err != nil {
				fmt.Println("failed to marshal")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("JSON not valid"))
			}

			rw.Write(end)
			return
		}

		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("Unauthorized"))

	}
}

func AddSongToBucket(db *sql.DB) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var t database.Playlist
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("Failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		res := database.CheckIfSongExits(db, t.Uri, t.PlaylistName)
		if !res {
			rw.WriteHeader(http.StatusAlreadyReported)
			end, err := json.Marshal("exists")
			if err != nil {
				fmt.Println("failed to marshal")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("JSON not valid"))
			}
			rw.Write(end)
			return
		}

		lid, err := database.Insert(db, t)
		if err != nil {
			fmt.Println("failed Insert")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}
		if lid == 0 {
			fmt.Println("Error Insert Id cannot be 0")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		end, err := json.Marshal(t)
		if err != nil {
			fmt.Println("failed to marshal")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		rw.Write(end)
	}
}

func (p *PlaylistCredentials) GetBucket(db *sql.DB) http.HandlerFunc {
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

		end, err := database.SelectBucket(db, t.User, t.PlaylistName, t.From, t.To)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Select has failed"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) GetPlaylist(db *sql.DB) http.HandlerFunc {
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
			fmt.Println(err)
		}

		end, err := database.SelectPlaylist(db, t.PlaylistName, t.From, t.To)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Select has failed"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) Delete(db *sql.DB) http.HandlerFunc {
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
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Delete has failed"))
		}
		anz, err := database.DeleteSong(db, t.ID, t.User)

		if err != nil || anz != 1 {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Delete has failed"))
		}

		end, err := json.Marshal("sucess")
		if err != nil {
			fmt.Println("failed to marshal")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}
		rw.Write(end)
	}
}

func (p *PlaylistCredentials) Dislike(db *sql.DB) http.HandlerFunc {
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

		del, errs := database.DislikeSong(db, t.ID, t.User, t.PlaylistName)

		if errs != nil {
			fmt.Println("failed to Update")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("update failed"))
		}

		ret := "not deleted"

		if del {
			ret = "deleted"
		}

		end, err := json.Marshal(ret)
		if err != nil {
			fmt.Println("failed to marshal")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		rw.Write(end)
	}
}
