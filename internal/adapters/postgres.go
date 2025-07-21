package adapters

import (
	"context"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasteli-dev/tvtrackr-go/internal/config"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB() (*PostgresDB, error) {
	url := config.PostgresURL
	if url == "" {
		url = os.Getenv("POSTGRES_URL")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
