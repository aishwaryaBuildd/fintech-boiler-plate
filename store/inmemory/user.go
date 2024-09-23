package inmemory

import (
	"errors"

	"github.com/aishwaryaBuildd/go_boiler_plate/store"
)

type InMemoryUserStore struct {
	users map[string]*store.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*store.User),
	}
}

func (s *InMemoryUserStore) CreateUser(user *store.User) error {
	if _, exists := s.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	s.users[user.Email] = user
	return nil
}

func (s *InMemoryUserStore) GetUserByEmail(email string) (*store.User, error) {
	user, exists := s.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
