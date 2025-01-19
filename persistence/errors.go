package persistence

import "errors"

// ErrNotFound is returned when a resource does not exist
// TODO: probably replace with a boolean as a return type for the repository
var ErrNotFound = errors.New("record not found")

// ErrIdKeyCollision is returned when trying to create a resource with the same ID.
var ErrIdKeyCollision = errors.New("id key collision")
