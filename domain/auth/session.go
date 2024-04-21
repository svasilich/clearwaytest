package auth

import "time"

// UserSession represent information about user session.
type UserSession struct {
	Token    Token
	OpenedAt time.Time
}
