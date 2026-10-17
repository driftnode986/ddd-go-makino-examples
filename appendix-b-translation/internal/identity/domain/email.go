// #region email
package domain

import (
	"fmt"
	"net/mail"
	"strings"
)

type Email struct {
	local  string
	domain string
}

func NewEmail(s string) (Email, error) {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return Email{}, fmt.Errorf("invalid email: %w", err)
	}
	at := strings.LastIndex(addr.Address, "@")
	return Email{
		local:  addr.Address[:at],
		domain: strings.ToLower(addr.Address[at+1:]),
	}, nil
}

// #endregion

func (e Email) String() string { return e.local + "@" + e.domain }
