package cwrepo

import "errors"

// ErrUserNotExists throw when can't find login and password pair in db.
var ErrUserNotExists = errors.New("user not exists")

// ErrNoOpenSessions throw when cat't find any open sessions for the user.
var ErrNoOpenSessions = errors.New("cat't find any open sessions for the user")

// ErrForbiddenAsset throw when cat't find asset with given credentials
var ErrForbiddenAsset = errors.New("cat't find asset with given credentials")
