package auth

import (
	"context"

	"github.com/svasilich/clearwaytest/domain/auth"
	dauth "github.com/svasilich/clearwaytest/domain/auth"
)

// DBUserAuthorizer сheck if the user exists in the database and create a new user session, or return an existing one.
type DBUserAuthorizer interface {
	Login(ctx context.Context, user string, passwordHash string) (dauth.UserSession, error)
}

// DBUserRetriever retrieve user by session if the session is active. Only the last user session is active.
type DBUserRetriever interface {
	GetUserBySession(ctx context.Context, token auth.Token) (int64, error)
}
