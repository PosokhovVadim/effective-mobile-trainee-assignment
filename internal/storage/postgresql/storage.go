package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"songs_lib/internal/model"
	"songs_lib/pkg/logger"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

type PostgresStorage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPostgresStorage(log *slog.Logger, path string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &PostgresStorage{
		db:  db,
		log: log,
	}, nil
}

func (s *PostgresStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (p *PostgresStorage) BeginTx() (*sql.Tx, error) {
	return p.db.Begin()
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil

}
func (s *PostgresStorage) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		s.log.Error("Transaction failed", logger.Err(err))
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) AddSong(song model.Song, verses []string) (uint, error) {
	var songID uint

	if err := s.WithTransaction(func(tx *sql.Tx) error {
		err := tx.QueryRow(
			`INSERT INTO songs (group_name, name, link, release_date, inserted_at) 
             VALUES ($1, $2, $3, $4, NOW()) 
             RETURNING id`,
			song.Group, song.Name, song.Link, song.ReleaseDate,
		).Scan(&songID)
		if err != nil {
			return err
		}

		for i, verse := range verses {
			_, err = tx.Exec(
				`INSERT INTO verses (song_id, verse_number, verse_text) 
                 VALUES ($1, $2, $3)`,
				songID, i+1, verse,
			)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return 0, err
	}
	s.log.Info("Song added successfully", slog.Int("song_id", int(songID)))

	return songID, nil
}

func (s *PostgresStorage) DeleteSong(songID uint) error {
	_, err := s.db.Exec(
		`DELETE FROM songs WHERE id = $1`,
		songID,
	)
	if err != nil {
		return err
	}
	s.log.Info("Song deleted successfully", slog.Int("song_id", int(songID)))

	return nil
}

func (s *PostgresStorage) GetLyrics(songID uint, limit, offset int) ([]model.Lyrics, error) {
	rows, err := s.db.Query(
		`SELECT song_id, verse_number, text FROM lyrics 
         WHERE song_id = $1 
         ORDER BY verse_number 
         LIMIT $2 OFFSET $3`,
		songID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lyrics []model.Lyrics
	for rows.Next() {
		var lyric model.Lyrics
		if err := rows.Scan(&lyric.SongID, &lyric.VerseNumber, &lyric.Text); err != nil {
			return nil, err
		}
		lyrics = append(lyrics, lyric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lyrics, nil
}
