// #region errors
package domain

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidStatus = errors.New("invalid order status")
)

// #endregion

// ErrInvalidOrder は注文の内容が不正なとき（明細が空など）のエラー。
var ErrInvalidOrder = errors.New("invalid order")
