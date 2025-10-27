package model

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrOrderCancelled = errors.New("order cancelled")
	ErrOrderPaid      = errors.New("order paid")
)
