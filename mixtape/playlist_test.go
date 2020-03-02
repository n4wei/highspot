package mixtape_test

import (
	"io"
	"log"

	mixtape_pkg "github.com/n4wei/highspot/mixtape"
	"github.com/n4wei/highspot/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Playlist Changes", func() {
	var (
		mixtape *models.Mixtape
		changes *models.Changes

		testOutput  io.Writer
		testMixtape *mixtape_pkg.Mixtape

		applyChangesErr error
	)

	BeforeEach(func() {
		mixtape = &models.Mixtape{
			Users: []models.User{
				{
					ID:   "user_1",
					Name: "test_user_1",
				},
				{
					ID:   "user_2",
					Name: "test_user_2",
				},
				{
					ID:   "user_3",
					Name: "test_user_3",
				},
			},
			Playlists: []models.Playlist{
				{
					ID:     "playlist_1",
					UserID: "user_1",
					SongIDs: []string{
						"song_1",
						"song_2",
					},
				},
				{
					ID:     "playlist_2",
					UserID: "user_3",
					SongIDs: []string{
						"song_3",
					},
				},
			},
			Songs: []models.Song{
				{
					ID:     "song_1",
					Artist: "some_artist",
					Title:  "test_song_1",
				},
				{
					ID:     "song_2",
					Artist: "some_other_artist",
					Title:  "test_song_2",
				},
				{
					ID:     "song_3",
					Artist: "another_artist",
					Title:  "test_song_3",
				},
			},
		}

		changes = &models.Changes{}

		testOutput = gbytes.NewBuffer()
		testMixtape = mixtape_pkg.New(mixtape, log.New(testOutput, "", 0))
	})

	JustBeforeEach(func() {
		applyChangesErr = testMixtape.ApplyChanges(changes)
	})

	Describe("addPlaylist", func() {
		Context("when the new playlist is missing ID", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id missing, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when the new playlist has an ID that already exists", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_1",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id playlist_1 already exists, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when the new playlist is missing UserID", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("user_id missing, from playlist_id playlist_x, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when the new playlist has an UserID not in mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_x",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("user_id user_x not in mixtape, from playlist_id playlist_x, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when the new playlist does not contain any songs", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{},
						},
					},
				}
			})

			It("should not add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id playlist_x does not contain any songs, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when one of the songs in the new playlist does not exist in the mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{"song_1", "song_x"},
						},
					},
				}
			})

			It("should not add non-existant song, add the playlist, output logs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("song_id song_x not in mixtape, from playlist_id playlist_x, skipping"))
				Expect(testOutput).To(gbytes.Say("added playlist_id playlist_x"))

				Expect(mixtape.Playlists).To(HaveLen(3))
				Expect(mixtape.Playlists[2]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_x",
					UserID:  "user_1",
					SongIDs: []string{"song_1"},
				}))
			})
		})

		Context("when the new playlist does not contain any songs from mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{"song_x", "song_y"},
						},
					},
				}
			})

			It("should not add the playlist, output logs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("song_id song_x not in mixtape, from playlist_id playlist_x, skipping"))
				Expect(testOutput).To(gbytes.Say("song_id song_y not in mixtape, from playlist_id playlist_x, skipping"))
				Expect(testOutput).To(gbytes.Say("playlist_id playlist_x does not contain any songs from mixtape, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when adding a valid playlist", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{"song_1", "song_2"},
						},
					},
				}
			})

			It("should add the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("added playlist_id playlist_x"))

				Expect(mixtape.Playlists).To(HaveLen(3))
				Expect(mixtape.Playlists[2]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_x",
					UserID:  "user_1",
					SongIDs: []string{"song_1", "song_2"},
				}))
			})
		})
	})

	Describe("removePlaylist", func() {
		Context("when playlist ID is missing", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "",
						},
					},
				}
			})

			It("should not remove any playlist, output a log, and continues", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id missing, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when playlist ID does not exist in mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "playlist_x",
						},
					},
				}
			})

			It("should not remove any playlist, output a log, and continues", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id playlist_x not found, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
			})
		})

		Context("when playlist ID exists in mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "playlist_2",
						},
					},
				}
			})

			It("should remove the playlist, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("removed playlist_id playlist_2"))

				Expect(mixtape.Playlists).To(HaveLen(1))
				Expect(mixtape.Playlists[0].ID).To(Equal("playlist_1"))
			})
		})
	})

	Describe("addSongsToPlaylist", func() {
		Context("when the playlist ID is missing", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "",
							SongIDs: []string{"song_1"},
						},
					},
				}
			})

			It("should not add any songs, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id missing, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[0].SongIDs).To(HaveLen(2))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(1))
			})
		})

		Context("when the playlist ID does not exist in mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add any songs, output a log, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("playlist_id playlist_x not found, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[0].SongIDs).To(HaveLen(2))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(1))
			})
		})

		Context("when the playlist does not contain any songs", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{},
						},
					},
				}
			})

			It("should not add any songs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).ToNot(gbytes.Say("added song_id"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(1))
			})
		})

		Context("when the playlist contains songs that are not in mixtape", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{"song_x"},
						},
					},
				}
			})

			It("should not add these songs, output logs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("song_id song_x not in mixtape, not added to playlist_id playlist_2, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(1))
			})
		})

		Context("when the playlist contains songs that are already in the playlist", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should not add these songs, output logs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("song_id song_3 already in playlist_id playlist_2, skipping"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(1))
			})
		})

		Context("when the playlist contains valid songs", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{"song_2", "song_1"},
						},
					},
				}
			})

			It("should add them to the playlist, output logs, and continue", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("added song_id song_2 to playlist_id playlist_2"))
				Expect(testOutput).To(gbytes.Say("added song_id song_1 to playlist_id playlist_2"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[1].ID).To(Equal("playlist_2"))
				Expect(mixtape.Playlists[1].SongIDs).To(HaveLen(3))
				Expect(mixtape.Playlists[1].SongIDs).To(ContainElements([]string{"song_1", "song_2", "song_3"}))
			})
		})

		Context("when user ID is provided", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_1",
							UserID:  "user_x",
							SongIDs: []string{"song_3"},
						},
					},
				}
			})

			It("should ignored user ID", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())
				Expect(testOutput).To(gbytes.Say("added song_id song_3 to playlist_id playlist_1"))

				Expect(mixtape.Playlists).To(HaveLen(2))
				Expect(mixtape.Playlists[0].ID).To(Equal("playlist_1"))
				Expect(mixtape.Playlists[0].UserID).To(Equal("user_1"))
			})
		})
	})

	Describe("composite changes", func() {
		Context("add, update", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							SongIDs: []string{"song_1", "song_3"},
						},
					},
				}
			})

			It("should do the right things", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())

				Expect(mixtape.Playlists).To(HaveLen(3))
				Expect(mixtape.Playlists[2]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_x",
					UserID:  "user_1",
					SongIDs: []string{"song_3", "song_1"},
				}))
			})
		})

		Context("update, add, remove", func() {
			BeforeEach(func() {
				changes.PlaylistChanges = []models.PlaylistChange{
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{"song_1", "song_x"},
						},
					},
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_2",
							SongIDs: []string{"song_1"},
						},
					},
					{
						ID: models.AddSongs,
						Playlist: models.Playlist{
							ID:      "playlist_1",
							SongIDs: []string{"song_2", "song_3"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_x",
							UserID:  "user_1",
							SongIDs: []string{"song_x", "song_1"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_y",
							UserID:  "user_2",
							SongIDs: []string{"song_2"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_z",
							UserID:  "user_1",
							SongIDs: []string{"song_z"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_a",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_b",
							UserID:  "user_b",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_c",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "playlist_c",
						},
					},
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "playlist_1",
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_d",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.Add,
						Playlist: models.Playlist{
							ID:      "playlist_e",
							UserID:  "user_1",
							SongIDs: []string{"song_3"},
						},
					},
					{
						ID: models.Remove,
						Playlist: models.Playlist{
							ID: "playlist_d",
						},
					},
				}
			})

			It("should do the right things", func() {
				Expect(applyChangesErr).ToNot(HaveOccurred())

				Expect(mixtape.Playlists).To(HaveLen(5))
				Expect(mixtape.Playlists[0]).To(BeEquivalentTo(models.Playlist{
					ID:     "playlist_a",
					UserID: "user_1",
					SongIDs: []string{
						"song_3",
					},
				}))
				Expect(mixtape.Playlists[1]).To(BeEquivalentTo(models.Playlist{
					ID:     "playlist_2",
					UserID: "user_3",
					SongIDs: []string{
						"song_3",
						"song_1",
					},
				}))
				Expect(mixtape.Playlists[2]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_x",
					UserID:  "user_1",
					SongIDs: []string{"song_1"},
				}))
				Expect(mixtape.Playlists[3]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_y",
					UserID:  "user_2",
					SongIDs: []string{"song_2"},
				}))
				Expect(mixtape.Playlists[4]).To(BeEquivalentTo(models.Playlist{
					ID:      "playlist_e",
					UserID:  "user_1",
					SongIDs: []string{"song_3"},
				}))
			})
		})
	})
})
