package domain

import "errors"

// ErrUserNotFound はリポジトリに会員が無いときのエラー。
var ErrUserNotFound = errors.New("user not found")
