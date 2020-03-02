package models

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Song struct {
	ID     string `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

type Playlist struct {
	ID      string   `json:"id"`
	UserID  string   `json:"user_id"`
	SongIDs []string `json:"song_ids"`
}

type Mixtape struct {
	Users     []User     `json:"users"`
	Playlists []Playlist `json:"playlists"`
	Songs     []Song     `json:"songs"`
}
