package postgresql

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

// Constructor for Storage
func New(connectionString string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	const op = "storage.postgres.Stop"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}