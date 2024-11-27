package dto

import (
	"songs_lib/internal/model"
	"time"
)

type CreateSongRequest struct {
	Group string `json:"group" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

type CreateSongResponse struct {
	ID          uint      `json:"id"`
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"release_date"`
	Link        string    `json:"link,omitempty"`
	Text        string    `json:"text,omitempty"`
}

type LyricsDTO struct {
	VerseNumber uint   `json:"verse_number,omitempty"`
	Text        string `json:"text,omitempty"`
}

type SongDTO struct {
	ID          uint        `json:"id,omitempty"`
	Group       string      `json:"group,omitempty"`
	Name        string      `json:"name,omitempty"`
	ReleaseDate time.Time   `json:"release_date,omitempty"`
	Link        string      `json:"link,omitempty"`
	InsertedAt  string      `json:"inserted_at,omitempty"`
	Lyrics      []LyricsDTO `json:"verses,omitempty"`
}

type LibraryDTO struct {
	Songs []SongDTO `json:"songs"`
}

func LyricsToDTO(lyrics model.Lyrics) LyricsDTO {
	return LyricsDTO{
		VerseNumber: lyrics.VerseNumber,
		Text:        lyrics.Text,
	}
}

func SongToDTO(song model.Song, lyrics []model.Lyrics) SongDTO {
	var lyricsDTO []LyricsDTO
	for _, lyric := range lyrics {
		lyricsDTO = append(lyricsDTO, LyricsToDTO(lyric))
	}

	return SongDTO{
		ID:          song.ID,
		Group:       song.Group,
		Name:        song.Name,
		ReleaseDate: song.ReleaseDate,
		Link:        song.Link,
		InsertedAt:  song.InsertedAt.Format("2006-01-02 15:04:05"),
		Lyrics:      lyricsDTO,
	}
}
