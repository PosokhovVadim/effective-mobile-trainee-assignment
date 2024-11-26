package postgresql

import "database/sql"

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(path string) (*PostgresStorage, error) {
	// connect to storage
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	//TODO: up migrations

	return &PostgresStorage{
		db: db,
	}, nil
}
