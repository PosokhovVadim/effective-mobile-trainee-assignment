package service

import (
	"log/slog"
	"songs_lib/internal/model"
	"songs_lib/internal/storage"
	"strings"
	"time"
)

type ISong interface {
	AddSong(group, name, link string, releaseDate time.Time, text string) (uint, error)
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

func (s *SongService) AddSong(group, name, link string, releaseDate time.Time, text string) (uint, error) {
	var verses []string
	if len(text) != 0 {
		verses = splitTextIntoVerses(text)
	}

	song := model.Song{
		Group:       group,
		Name:        name,
		ReleaseDate: releaseDate,
		Link:        link,
	}

	songID, err := s.s.AddSong(song, verses)
	if err != nil {
		return 0, err
	}
	return songID, nil
}

func splitTextIntoVerses(text string) []string {
	verses := strings.Split(text, "\n\n")

	for i, verse := range verses {
		verses[i] = strings.TrimSpace(verse)
	}

	filtered := make([]string, 0)
	for _, verse := range verses {
		if verse != "" {
			filtered = append(filtered, verse)
		}
	}
	return filtered
}
