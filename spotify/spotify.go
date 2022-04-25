package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// NewSpotifyInstance creates a new Instance with AcessToken
func NewSpotifyInstance(flag *bool) SpotifyCredentials {
	spotifyInstance := SpotifyCredentials{
		accessToken: "",
		headerFlag:  flag,
	}

	err := spotifyInstance.getAccessToken()
	if err != nil {
		fmt.Println("FAIL GETTING ACCESS TOKEN!")
	}

	fmt.Println("new Spotify Instance created")

	return spotifyInstance
}

// Get Request to Spotify API to get Access Token
func (s *SpotifyCredentials) getAccessToken() error {
	resp, err := http.Get("https://open.spotify.com/get_access_token?reason=transport&productType=web_player")
	if err != nil {
		return err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data := Acess{}

	err2 := json.Unmarshal(rbody, &data)
	if err2 != nil {
		return err2
	}

	s.accessToken = data.AccessToken

	return nil
}

func (s *SpotifyCredentials) getSong(value string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?q="+url.QueryEscape(value)+"&type=track&market=US&limit=10&offset=0", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	if err != nil {
		fmt.Println("failed to send request")
		return nil, 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed to get response")
		return nil, 0, err
	}

	ok, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read response")
		return nil, 0, err
	}

	defer resp.Body.Close()
	data := SpotifyItem{}
	err2 := json.Unmarshal(ok, &data)
	if err2 != nil {
		fmt.Println("Fail in Unmarshal JSON")
		return nil, 0, err2
	}

	result := data.Tracks.Items

	end, err := json.Marshal(result)
	if err != nil {
		fmt.Println("FAILED MARSHAL JSON")
		return nil, 0, err
	}

	return end, resp.StatusCode, nil
}
