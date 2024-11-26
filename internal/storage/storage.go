package storage

import (
	"songs_lib/internal/model"
)

type Storage interface {
	AddSong(song model.Song, verses []string) (uint, error)
	DeleteSong(songID uint) error
	GetLyrics(songID uint, limit, offset int) ([]model.Lyrics, error)
	GetAllSongs() ([]model.Song, error)
	GetAllLyrics() ([]model.Lyrics, error)
	UpdateSong(song model.Song) error
	UpdateLyrics(lyrics model.Lyrics) error
}
