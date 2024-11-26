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
	Text        string    `json:"text,omitempty"`
}

type LyricsInfo struct {
	VerseNumber uint   `json:"verse_number"`
	Text        string `json:"text"`
}

type GetLyricsResponse struct {
	SongName string       `json:"song_name"`
	Group    string       `json:"group"`
	Verses   []LyricsInfo `json:"verses"`
}
