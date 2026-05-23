// shared/dto/lesson_view.go
package dto

import "time"

// LessonView は Catalog BC が他 BC 向けに公開する「公表された言語」DTO です。
// Booking BC は内部の domain.Lesson を import せず、この LessonView だけを参照します。
type LessonView struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	PriceJPY  int       `json:"price_jpy"`
	StartsAt  time.Time `json:"starts_at"`
	TrainerID string    `json:"trainer_id"`
}
