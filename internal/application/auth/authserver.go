package auth

import (
	"context"
	"errors"
	"fmt"

	dauth "github.com/svasilich/clearwaytest/domain/auth"
	"github.com/svasilich/clearwaytest/domain/common"
)

// Auth work with users credentials.
type Auth struct {
	authorizerRepo DBUserAuthorizer
	hasherFunc     Hasher
}

// NewAuth create instance of Auth.
func NewAuth(authorizerRepo DBUserAuthorizer, hasher Hasher) *Auth {
	return &Auth{
		authorizerRepo: authorizerRepo,
		hasherFunc:     hasher,
	}
}

// Login is handle user log in.
func (a *Auth) Login(ctx context.Context, user string, password string) (dauth.Token, error) {
	hash, err := a.hasherFunc(password)
	if err != nil {
		return dauth.Token(""), fmt.Errorf("can't make hash for pair (%s:%s): %w", user, password, err)
	}
	sessionString, err := a.authorizerRepo.Login(ctx, user, hash)
	if err != nil {
		if errors.Is(err, common.ErrUnauthorized) {
			return dauth.Token(""), err
		}

		return dauth.Token(""), fmt.Errorf("unexpected error when %s logging in: %w", user, err)
	}
	return sessionString.Token, nil
}
