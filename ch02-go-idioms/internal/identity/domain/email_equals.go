package domain

// Equals は値が等しいかを返す。2.6 節の User.ChangeEmail が使う
// (値オブジェクトの等価性は第3章で詳しく扱う)。
func (e Email) Equals(other Email) bool { return e.value == other.value }
