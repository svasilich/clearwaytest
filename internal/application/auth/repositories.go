package auth

import (
	"context"

	dauth "github.com/svasilich/clearwaytest/domain/auth"
)

// DBUserAuthorizer сheck if the user exists in the database and create a new user session, or return an existing one.
type DBUserAuthorizer interface {
	Login(ctx context.Context, user string, passwordHash string) (dauth.UserSession, error)
}

// DBUserSessionRetriever retrieves the last opened session for a usee.
type DBUserSessionRetriever interface {
	GetLastSession(ctx context.Context, userID int64) (dauth.UserSession, error)
}
