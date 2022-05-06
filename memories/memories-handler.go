package memories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"playlisttogether/backend/database"
	"playlisttogether/backend/playlist"
	"playlisttogether/backend/utils"
)

type ImagesResponse struct {
	Base64Code string `json:"base64Code"`
}

func GetMemories(db *sql.DB, imgpath string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		(rw).Header().Set("Access-Control-Allow-Origin", "*")
		(rw).Header().Set("Access-Control-Allow-Methods", "POST")
		(rw).Header().Set("Access-Control-Allow-Headers", "*")

		if req.Method == "OPTIONS" {
			rw.WriteHeader(http.StatusOK)
			return
		}
		decoder := json.NewDecoder(req.Body)
		var t playlist.Auth
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println("failed to decode")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}
		md5folder := utils.GetMD5Hash(t.PlaylistName)

		paths, err2 := database.SelectPictures(db, md5folder, t.From, t.To)
		if err2 != nil {
			fmt.Println("failed Selecting pictures")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("JSON not valid"))
		}

		allImages := getFilesPaths(paths)

		rw.Write(allImages)
	}

}
