package http

type RegisterUserRequest struct {
	Email string `json:"email"` // 公開フィールド、検証は handler で
}
