package core

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type userServiceMem struct {
	mu             sync.Mutex
	users          map[string]*User           // key: username
	seriesFollowed map[string]map[string]bool // userID -> set of seriesID
}

func NewUserServiceMem() UserService {
	return &userServiceMem{
		users:          make(map[string]*User),
		seriesFollowed: make(map[string]map[string]bool),
	}
}

func (s *userServiceMem) Register(ctx context.Context, username, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[username]; exists {
		return nil, errors.New("username already exists")
	}
	user := &User{
		ID:       uuid.NewString(),
		Username: username,
		Password: password, // In real code, hash this!
	}
	s.users[username] = user
	return user, nil
}

func (s *userServiceMem) Login(ctx context.Context, username, password string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, exists := s.users[username]
	if !exists || user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *userServiceMem) FollowSeries(ctx context.Context, userID, seriesID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.seriesFollowed[userID]; !ok {
		s.seriesFollowed[userID] = make(map[string]bool)
	}
	s.seriesFollowed[userID][seriesID] = true
	return nil
}

func (s *userServiceMem) ListFollowedSeries(ctx context.Context, userID string) ([]Series, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	seriesIDs, ok := s.seriesFollowed[userID]
	if !ok {
		return nil, nil
	}
	var result []Series
	for id := range seriesIDs {
		result = append(result, Series{ID: id, Title: ""}) // Title unknown in memory implementation
	}
	return result, nil
}
