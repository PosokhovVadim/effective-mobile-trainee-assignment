package storage

import (
	"database/sql"
	"songs_lib/internal/model"
)

type Storage interface {
	AddSong(song model.Song, verses []string) (uint, error)
	BeginTx() (*sql.Tx, error)
}
