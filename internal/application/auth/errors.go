package auth

import "errors"

// ErrUnauthorized throw when can't find login and password pair in db.
var ErrUnauthorized = errors.New("user not authorized")

// ErrSessionNotActive throw when given session is not active.
var ErrSessionNotActive = errors.New("session is not active")

// ErrSessionNotExists throw when given session is not exists.
var ErrSessionNotExists = errors.New("session not exists")
