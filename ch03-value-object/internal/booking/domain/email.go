// #region email-type
package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// Email は値オブジェクト。フィールドは非公開で、コンストラクタ経由でのみ生成される。
type Email struct {
	local  string
	domain string
}

// ErrInvalidEmail は不正なメール形式を表すセンチネルエラー。
var ErrInvalidEmail = errors.New("invalid email")

// #endregion

// #region new-email
// NewEmail は RFC 5322 に従ってメールアドレスを検証し、Email を返す。
func NewEmail(s string) (Email, error) {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return Email{}, fmt.Errorf("%w: %v", ErrInvalidEmail, err)
	}
	at := strings.LastIndex(addr.Address, "@")
	if at < 0 {
		return Email{}, fmt.Errorf("%w: missing @", ErrInvalidEmail)
	}
	return Email{
		local:  addr.Address[:at],
		domain: strings.ToLower(addr.Address[at+1:]),
	}, nil
}

// #endregion

// #region email-methods
// String は表示用文字列を返す (fmt.Stringer)。
func (e Email) String() string {
	if e.local == "" && e.domain == "" {
		return ""
	}
	return e.local + "@" + e.domain
}

// Equals は等価性を判定する。ドメイン部は大文字小文字を区別しない。
func (e Email) Equals(other Email) bool {
	return e.local == other.local && e.domain == other.domain
}

// Domain はメールアドレスのドメイン部を返す。読み取り専用。
func (e Email) Domain() string {
	return e.domain
}

// #endregion
