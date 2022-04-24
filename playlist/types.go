package playlist

type Auth struct {
	User         string `json:"user"`
	Token        string `json:"token"`
	ID           int    `json:"id"`
	PlaylistName string `json:"playlistname"`
	PlaylistUrl  string `json:"playlisturl"`
}
