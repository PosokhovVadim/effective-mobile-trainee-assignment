package dto

import "time"

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
	InsertedAt  time.Time `json:"inserted_at"`
}

type GetLyricsResponse struct {
	SongID      uint   `json:"song_id"`
	VerseNumber uint   `json:"verse_number"`
	Text        string `json:"text"`
}
