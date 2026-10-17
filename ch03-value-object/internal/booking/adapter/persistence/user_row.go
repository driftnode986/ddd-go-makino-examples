package persistence

type UserRow struct {
	ID    string `db:"id"`
	Email string `db:"email"` // 公開フィールド、Persistence Model
}
