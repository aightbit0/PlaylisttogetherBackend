package database

type PlaylistUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Active   bool   `json:"active"`
}

type Playlist struct {
	ID           int    `json:"id"`
	User         string `json:"name"`
	Artist       string `json:"artist"`
	Url          string `json:"url"`
	Picture      string `json:"picture"`
	Uri          string `json:"uri"`
	Songname     string `json:"songname"`
	Playlist     string `json:"playlust"`
	Dislike      int    `json:"dislike"`
	Disliker     string `json:"disliker"`
	PlaylistName string `json:"playlistname"`

	Token string `json:"token"`
	//Dislikes     string `json:"dislikes"`
}

type Playlists struct {
	ID           int    `json:"id"`
	User         string `json:"user"`
	PlaylistName string `json:"playlistname"`
	PlaylistUrl  string `json:"playlisturl"`
	PlaylistID   string `json:"playlistid"`
	Amount       int    `json:"amount"`
	Creator      string `json:"creator"`

	Users []Users `json:"users"`
}

type Users struct {
	Label string `json:"label"`
}

type Bucket struct {
	Data   []Playlist `json:"data"`
	Amount int        `json:"amount"`
}
