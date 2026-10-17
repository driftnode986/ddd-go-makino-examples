package domain

import (
	"errors"
	"fmt"
	"net/mail"
)

// Email は値オブジェクト。フィールドは非公開で、コンストラクタ経由でのみ生成される
type Email struct {
	value string
}

var ErrInvalidEmail = errors.New("invalid email")

func NewEmail(s string) (Email, error) {
	if _, err := mail.ParseAddress(s); err != nil {
		return Email{}, fmt.Errorf("%w: %v", ErrInvalidEmail, err)
	}
	return Email{value: s}, nil
}

func (e Email) String() string { return e.value }
