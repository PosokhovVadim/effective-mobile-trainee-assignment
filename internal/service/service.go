package service

import (
	"log/slog"
	"songs_lib/internal/storage"
	"time"
)

type ISong interface {
	AddSong(group, name, link string, releaseDate time.Time, text string) error
}

type SongService struct {
	s   storage.Storage
	log *slog.Logger
}

func NewSongService(log *slog.Logger, s storage.Storage) *SongService {
	return &SongService{
		log: log,
		s:   s,
	}
}

func (s *SongService) AddSong(group, name, link string, releaseDate time.Time, text string) error {
	return nil
}
