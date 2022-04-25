package spotify

import (
	"fmt"
	"net/http"
)

// Searcher
func (s *SpotifyCredentials) Searcher(rw http.ResponseWriter, req *http.Request) {
	if *s.headerFlag {
		(rw).Header().Set("Access-Control-Allow-Origin", "*")
		(rw).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		(rw).Header().Set("Access-Control-Allow-Headers", "*")
	}
	keys, ok := req.URL.Query()["v"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("Url Param 'key' is missing")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Url Param 'key' is missing"))
		return
	}
	key := keys[0]
	fin, statuscode, err := s.getSong(string(key))

	if err != nil {
		fmt.Println("FAILED GET SONG")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("FAILED GET SONG"))
		return
	}

	// renew Access Token
	if statuscode == 401 {

		//fmt.Println("REQUESTING NEW ACESS TOKEN")
		err := s.getAccessToken()

		if err != nil {
			fmt.Println("FAIL TO RENEW ACCESS TOKEN")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("FAIL TO RENEW ACCESS TOKEN"))
			return
		}

		fin, statuscode, err = s.getSong(string(key))

		if err != nil || statuscode == 401 {
			fmt.Println("FAIL TO RENEW ACCESS TOKEN")
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Println("FAIL TO RENEW ACCESS TOKEN")
			return
		}
	}

	rw.Write(fin)
}
