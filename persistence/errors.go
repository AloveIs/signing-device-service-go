package persistence

import "errors"

// TODO: probably replace with a boolean as a return type for the repository
var ErrNotFound = errors.New("record not found")
