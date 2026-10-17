package application

import "fmt"

// ApplicationError はユースケース層が返すエラー。
// Presentation 層が HTTP ステータスにマッピングするときの
// 判断材料を持つ。
type ApplicationError struct {
	Code    string // "INSUFFICIENT_STOCK", "ORDER_NOT_FOUND", etc.
	Message string // ユーザー向けメッセージ
	Cause   error  // 原因 (ドメインエラー)
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ApplicationError) Unwrap() error {
	return e.Cause
}
