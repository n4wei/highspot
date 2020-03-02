package models

const (
	Add      PlaylistChangeID = "add"
	Remove   PlaylistChangeID = "remove"
	AddSongs PlaylistChangeID = "add_songs"
)

type PlaylistChangeID string

type PlaylistChange struct {
	ID       PlaylistChangeID `json:"id"`
	Playlist Playlist         `json:"playlist"`
}

type Changes struct {
	PlaylistChanges []PlaylistChange `json:"playlist_changes"`
}
