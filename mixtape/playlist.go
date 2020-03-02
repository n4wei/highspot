package mixtape

import "github.com/n4wei/highspot/models"

// These methods intentionally always return nil (instead of an error)
// because the UX design I chose is to skip invalid changes, log them,
// and keep going.

// This method adds a new playlist to the playlist array.
// It appends the new playlist to the end of the playlist array. A hash
// map is used to store the index of the playlist in the array for constant
// time access. If the new playlist has a user id that is not in mixtape,
// the playlist is not added. Only songs that exist in the mixtape are
// added with the new playlist. A playlist with 0 valid songs is not added.

// See tests in playlist_test.go for all invalid cases.

// runtime: O(s), s is the number of songs in the added playlist
// space: O(s), creates a map to store which songs belong to this playlist for
// constant time access
func (m *Mixtape) addPlaylist(playlist models.Playlist) error {
	m.logger.SetPrefix("[AddPlaylist] ")

	id := playlist.ID
	if id == "" {
		m.logger.Printf("playlist_id missing, skipping\n")
		return nil
	}
	if _, exist := m.lookup.playlists[id]; exist {
		m.logger.Printf("playlist_id %s already exists, skipping\n", id)
		return nil
	}
	if playlist.UserID == "" {
		m.logger.Printf("user_id missing, from playlist_id %s, skipping\n", id)
		return nil
	}
	if _, exist := m.lookup.users[playlist.UserID]; !exist {
		m.logger.Printf("user_id %s not in mixtape, from playlist_id %s, skipping\n", playlist.UserID, id)
		return nil
	}
	if len(playlist.SongIDs) == 0 {
		m.logger.Printf("playlist_id %s does not contain any songs, skipping\n", id)
		return nil
	}

	validSongIDs := []string{}
	for _, songID := range playlist.SongIDs {
		if _, exist := m.lookup.songs[songID]; exist {
			if m.lookup.playlistSongs[id] == nil {
				m.lookup.playlistSongs[id] = map[string]bool{}
			}
			m.lookup.playlistSongs[id][songID] = true
			validSongIDs = append(validSongIDs, songID)
		} else {
			m.logger.Printf("song_id %s not in mixtape, from playlist_id %s, skipping\n", songID, id)
		}
	}

	if len(validSongIDs) == 0 {
		m.logger.Printf("playlist_id %s does not contain any songs from mixtape, skipping\n", id)
		return nil
	}

	playlist.SongIDs = validSongIDs
	m.mixtape.Playlists = append(m.mixtape.Playlists, playlist)
	m.lookup.playlists[id] = len(m.mixtape.Playlists) - 1

	m.logger.Printf("added playlist_id %s\n", id)
	return nil
}

// This method removes an existing playlist.
// It does so by swapping the playlist we want to remove in the playlist array
// with the last element, and reslicing the array to reduce array length by 1.
// This is fast and requires constant time, however, it does not preserve the ordering
// of the playlist array elements. I decided to make this tradeoff because I assumed
// there would be many changes, and optimized for a fast remove operation.

// See tests in playlist_test.go for all invalid cases.

// runtime: O(1)
// space: no additional space
func (m *Mixtape) removePlaylist(playlist models.Playlist) error {
	m.logger.SetPrefix("[RemovePlaylist] ")

	id := playlist.ID
	if id == "" {
		m.logger.Printf("playlist_id missing, skipping\n")
		return nil
	}

	i, exist := m.lookup.playlists[id]
	if !exist {
		m.logger.Printf("playlist_id %s not found, skipping\n", id)
		return nil
	}

	playlists := m.mixtape.Playlists
	l := len(playlists)
	if i != l-1 {
		playlists[i], playlists[l-1] = playlists[l-1], playlists[i]
		m.lookup.playlists[playlists[i].ID] = i
	}

	m.mixtape.Playlists = playlists[:l-1]
	delete(m.lookup.playlists, id)
	delete(m.lookup.playlistSongs, id)

	m.logger.Printf("removed playlist_id %s\n", id)
	return nil
}

// This method adds one or more existing songs in mixtape to an existing
// playlist in mixtape.
// Songs are appended to the end of the playlist's list of songs. If a song
// does not exist in the mixtape or is already in the playlist, it is not added.

// See tests in playlist_test.go for all invalid cases.

// runtime: O(s), s is the number of songs being added
// space: creates a new entry for each valid song in the lookup hash map
// for songs belonging to this playlist
func (m *Mixtape) addSongsToPlaylist(playlist models.Playlist) error {
	m.logger.SetPrefix("[AddSongToPlaylist] ")

	id := playlist.ID
	if id == "" {
		m.logger.Printf("playlist_id missing, skipping\n")
		return nil
	}

	i, exist := m.lookup.playlists[id]
	if !exist {
		m.logger.Printf("playlist_id %s not found, skipping\n", id)
		return nil
	}

	for _, songID := range playlist.SongIDs {
		if _, exist = m.lookup.songs[songID]; !exist {
			m.logger.Printf("song_id %s not in mixtape, not added to playlist_id %s, skipping\n", songID, id)
			continue
		}
		if _, exist = m.lookup.playlistSongs[id][songID]; exist {
			m.logger.Printf("song_id %s already in playlist_id %s, skipping\n", songID, id)
			continue
		}

		m.mixtape.Playlists[i].SongIDs = append(m.mixtape.Playlists[i].SongIDs, songID)
		m.lookup.playlistSongs[id][songID] = true
		m.logger.Printf("added song_id %s to playlist_id %s\n", songID, id)
	}

	return nil
}
