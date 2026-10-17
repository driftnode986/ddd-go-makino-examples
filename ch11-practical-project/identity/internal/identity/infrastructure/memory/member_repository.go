package memory

import (
	"context"
	"sync"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/domain"
)

type MemberRepository struct {
	mu      sync.RWMutex
	members map[domain.MemberID]*domain.Member
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		members: map[domain.MemberID]*domain.Member{},
	}
}

func (r *MemberRepository) Save(
	_ context.Context,
	m *domain.Member,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.members[m.ID()] = m
	return nil
}

func (r *MemberRepository) FindByID(_ context.Context,
	id domain.MemberID) (*domain.Member, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.members[id]
	if !ok {
		return nil, domain.ErrMemberNotFound
	}
	return m, nil
}
