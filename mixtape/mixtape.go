package mixtape

import (
	"github.com/n4wei/highspot/models"
	"github.com/n4wei/highspot/util"
)

// Hash maps used for constant time lookup
type lookup struct {
	// map of user id to corresponding index in Mixtape.Users array
	users map[string]int
	// map of song id to corresponding index in Mixtape.Songs array
	songs map[string]int
	// map of playlist id to corresponding index in Mixtape.Playlists
	playlists map[string]int
	// map of playlist id to a second map of song ids belonging to this playlist
	playlistSongs map[string]map[string]bool
}

type Mixtape struct {
	mixtape *models.Mixtape
	lookup  *lookup

	logger util.Logger
}

func New(mixtape *models.Mixtape, logger util.Logger) *Mixtape {
	mt := &Mixtape{
		mixtape: mixtape,
		logger:  logger,
	}
	mt.buildLookup()
	return mt
}

// This method builds the lookup hash maps for the mixtape.
// It is created once in the constructor of the Mixtape object.
// runtime: O(u + s + p + p*ps)
// space: O(u + s + p + p*ps) since we create a hash map for each
// - u is the number of users
// - s is the number of songs
// - p is the number of playlists
// - ps is the most number of songs in any playlist
func (m *Mixtape) buildLookup() {
	lookup := &lookup{
		users:         map[string]int{},
		songs:         map[string]int{},
		playlists:     map[string]int{},
		playlistSongs: map[string]map[string]bool{},
	}

	for i, user := range m.mixtape.Users {
		lookup.users[user.ID] = i
	}
	for i, song := range m.mixtape.Songs {
		lookup.songs[song.ID] = i
	}
	for i, playlist := range m.mixtape.Playlists {
		lookup.playlists[playlist.ID] = i
		for _, songID := range playlist.SongIDs {
			if lookup.playlistSongs[playlist.ID] == nil {
				lookup.playlistSongs[playlist.ID] = map[string]bool{}
			}
			lookup.playlistSongs[playlist.ID][songID] = true
		}
	}

	m.lookup = lookup
}

// This method takes all the changes and applies them to mixtape in the
// order they were provided in the changes JSON file (order in an array).
func (m *Mixtape) ApplyChanges(changes *models.Changes) error {
	// I chose to go with the UX design of skipping invalid changes,
	// logging them, and keep applying further changes.
	// It's straightforward to make any of these methods return intentional
	// types of errors to stop applying changes partially through.
	var err error
	for _, change := range changes.PlaylistChanges {
		switch change.ID {
		case models.Add:
			err = m.addPlaylist(change.Playlist)
			if err != nil {
				return err
			}
		case models.Remove:
			err = m.removePlaylist(change.Playlist)
			if err != nil {
				return err
			}
		case models.AddSongs:
			err = m.addSongsToPlaylist(change.Playlist)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
