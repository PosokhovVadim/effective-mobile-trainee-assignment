package service

import (
	"log/slog"
	"songs_lib/internal/model"
	"songs_lib/internal/storage"
	"songs_lib/pkg/logger"
	"strings"
	"time"
)

type ISong interface {
	AddSong(group, name, link string, releaseDate time.Time, text string) (uint, error)
	DeleteSong(songID uint) error
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
		s.log.Error("Failed to add song", logger.Err(err))
		return 0, err
	}
	return songID, nil
}

func (s *SongService) DeleteSong(songID uint) error {
	if err := s.s.DeleteSong(songID); err != nil {
		s.log.Error("Failed to delete song", slog.Int("song_id", int(songID)), logger.Err(err))
		return err
	}
	return nil
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
