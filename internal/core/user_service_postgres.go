package core

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userServicePostgres struct {
	db *pgxpool.Pool
}

func NewUserServicePostgres(db *pgxpool.Pool) UserService {
	return &userServicePostgres{db: db}
}

func (s *userServicePostgres) Register(ctx context.Context, username, password string) (*User, error) {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	var id string
	err = s.db.QueryRow(ctx, `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`, username, string(hashed)).Scan(&id)
	// Check for unique violation (duplicate username)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		return nil, errors.New("username already exists")
	}
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Username: username, Password: ""}, nil // Do not return hash
}

func (s *userServicePostgres) Login(ctx context.Context, username, password string) (*User, error) {
	var user User
	err := s.db.QueryRow(ctx, `SELECT id, username, password FROM users WHERE username=$1`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	// Compare hash
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	user.Password = "" // Do not return hash
	return &user, nil
}

func (s *userServicePostgres) FollowSeries(ctx context.Context, userID, seriesID string) error {
	_, err := s.db.Exec(ctx, `INSERT INTO user_series (user_id, series_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, seriesID)
	return err
}

func (s *userServicePostgres) ListFollowedSeries(ctx context.Context, userID string) ([]Series, error) {
	rows, err := s.db.Query(ctx, `SELECT series_id FROM user_series WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Series
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, Series{ID: id, Title: ""})
	}
	return result, nil
}
