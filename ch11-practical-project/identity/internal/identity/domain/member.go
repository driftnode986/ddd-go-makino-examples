package domain

import (
	"errors"
	"net/mail"

	"github.com/google/uuid"
)

type MemberID struct{ value string }

func NewMemberID() MemberID { return MemberID{value: uuid.NewString()} }

func MemberIDFrom(s string) (MemberID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return MemberID{}, errors.New("invalid member id")
	}
	return MemberID{value: s}, nil
}

func (id MemberID) String() string { return id.value }

type Email struct{ value string }

func NewEmail(s string) (Email, error) {
	if _, err := mail.ParseAddress(s); err != nil {
		return Email{}, errors.New("invalid email")
	}
	return Email{value: s}, nil
}

func (e Email) String() string { return e.value }

type Member struct {
	id    MemberID
	email Email
	name  string
}

func NewMember(email Email, name string) (*Member, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	return &Member{id: NewMemberID(), email: email, name: name}, nil
}

func (m *Member) ID() MemberID { return m.id }
func (m *Member) Email() Email { return m.email }
func (m *Member) Name() string { return m.name }
