package store

import "errors"

// ErrKeyNotFound is returned when a requested key does not exist.
var ErrKeyNotFound = errors.New("key not found")

// Store defines the standard operations for our KV database.
type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}