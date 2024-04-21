package cwrepo

import "errors"

// ErrUserNotExists throw when can't find login and password pair in db.
var ErrUserNotExists = errors.New("user not exists")

// ErrNoOpenSessions не найдено открытых сессий для выбранного пользователя.
var ErrNoOpenSessions = errors.New("cat't find any open sessions for the user")
