package core

import "context"

type UserService interface {
	Register(ctx context.Context, username, password string) (*User, error)
	Login(ctx context.Context, username, password string) (*User, error)
	FollowSeries(ctx context.Context, userID, seriesID string) error
	ListFollowedSeries(ctx context.Context, userID string) ([]Series, error)
}
