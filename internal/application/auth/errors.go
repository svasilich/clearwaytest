package auth

import "errors"

// ErrUnauthorized throw when can't find login and password pair in db.
var ErrUnauthorized = errors.New("user not authorized")
