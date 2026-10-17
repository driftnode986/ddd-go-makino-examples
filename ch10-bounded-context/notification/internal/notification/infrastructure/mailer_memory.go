package infrastructure

import (
	"context"
	"sync"
)

type MailRecord struct {
	CustomerID string
	Body       string
}

type MemoryMailer struct {
	mu      sync.Mutex
	records []MailRecord
}

func NewMemoryMailer() *MemoryMailer {
	return &MemoryMailer{}
}

func (m *MemoryMailer) Send(
	ctx context.Context,
	customerID, body string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(
		m.records,
		MailRecord{CustomerID: customerID, Body: body},
	)
	return nil
}

func (m *MemoryMailer) SentCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.records)
}

func (m *MemoryMailer) Records() []MailRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]MailRecord, len(m.records))
	copy(out, m.records)
	return out
}
