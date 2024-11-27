package service

import (
	"log/slog"
	"songs_lib/internal/dto"
	"songs_lib/internal/model"
	"songs_lib/internal/storage"
	"songs_lib/pkg/logger"
	"strconv"
	"strings"
	"time"
)

type ISong interface {
	AddSong(group, name, link string, releaseDate time.Time, text string) (uint, error)
	DeleteSong(songID uint) error
	GetLyrics(songID uint, limit, offset string) (*dto.SongDTO, error)
	GetLibrary(filters map[string]string, limit, offset string) (*dto.LibraryDTO, error)
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

func (s *SongService) AddSong(
	group,
	name,
	link string,
	releaseDate time.Time,
	text string,
) (uint, error) {
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

func (s *SongService) GetLyrics(songID uint, limit, offset string) (*dto.SongDTO, error) {
	limitInt, offsetInt := getLimitAndOffset(limit, offset)
	song, err := s.s.GetSong(songID)
	if err != nil {
		s.log.Error("Failed to get song", logger.Err(err))
		return nil, err
	}

	lyrics, err := s.s.GetLyrics(songID, limitInt, offsetInt)
	if err != nil {
		s.log.Error("Failed to get lyrics", logger.Err(err))
		return nil, err
	}

	songDTO := &dto.SongDTO{
		Group: song.Group,
		Name:  song.Name,
	}

	for _, lyric := range lyrics {
		songDTO.Lyrics = append(songDTO.Lyrics, dto.LyricsDTO{
			VerseNumber: lyric.VerseNumber,
			Text:        lyric.Text,
		})
	}

	return songDTO, nil
}

func (s *SongService) GetLibrary(
	filters map[string]string,
	limit,
	offset string,
) (*dto.LibraryDTO, error) {
	limitInt, offsetInt := getLimitAndOffset(limit, offset)
	songsDTO := make([]dto.SongDTO, 0)

	songs, err := s.s.GetAllSongs(filters, limitInt, offsetInt)
	if err != nil {
		s.log.Error("Failed to get all songs", logger.Err(err))
		return nil, err
	}

	for _, song := range songs {
		lyrics, err := s.s.GetAllSongLyrics(song.ID)
		if err != nil {
			s.log.Error("Failed to get all song lyrics", logger.Err(err))
			return nil, err
		}

		songsDTO = append(songsDTO, dto.SongToDTO(song, lyrics))
	}

	return &dto.LibraryDTO{Songs: songsDTO}, nil
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

func getLimitAndOffset(limit, offset string) (int, int) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 10, 0
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return 10, 0
	}
	return limitInt, offsetInt
}
