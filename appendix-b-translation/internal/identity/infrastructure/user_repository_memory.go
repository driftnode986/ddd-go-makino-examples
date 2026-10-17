package infrastructure

import (
	"context"
	"errors"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
)

// MemoryUserRepository はテストで使う in-memory のリポジトリ。
type MemoryUserRepository struct {
	mu    sync.Mutex
	users map[domain.UserID]*domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{users: map[domain.UserID]*domain.User{}}
}

// ErrDuplicateUser は同じ ID の会員を 2 回 Add したときのエラー。
var ErrDuplicateUser = errors.New("user already exists")

func (r *MemoryUserRepository) Add(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.ID()]; ok {
		return ErrDuplicateUser
	}
	r.users[u.ID()] = u
	return nil
}

func (r *MemoryUserRepository) Save(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.ID()]; !ok {
		return domain.ErrUserNotFound
	}
	r.users[u.ID()] = u
	return nil
}

func (r *MemoryUserRepository) FindByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}

func (r *MemoryUserRepository) FindByEmail(_ context.Context, email domain.Email) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

var _ domain.UserRepository = (*MemoryUserRepository)(nil)
