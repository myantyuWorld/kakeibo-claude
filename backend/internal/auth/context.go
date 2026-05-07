package auth

import (
	"context"

	"backend/pkg/errors"
)

type contextKey struct{}

var userKey = contextKey{}

type AuthUser struct {
	ID         int64
	CognitoSub string
}

func GetUser(ctx context.Context) (AuthUser, error) {
	user, ok := ctx.Value(userKey).(AuthUser)
	if !ok {
		return AuthUser{}, errors.NewUnauthorized()
	}
	return user, nil
}
