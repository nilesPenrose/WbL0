package domain

import "errors"

var (
	ErrGetFromStorage = errors.New("can't get from Storage")
	ErrGetFromCache   = errors.New("can't get from OrderCache")
)

type ContextError struct {
	Error error
}
