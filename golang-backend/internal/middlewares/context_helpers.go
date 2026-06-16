package middlewares

import (
	"context"

	"github.com/aportela/doneo/internal/domain"
)

type contextKey string

type ContextUser struct {
	domain.UserBase
	SkipAuthorization bool
}

const contextUserKey contextKey = "user"

func GetContextUser(ctx context.Context) (ContextUser, bool) {
	user, ok := ctx.Value(contextUserKey).(ContextUser)
	return user, ok
}

func SetContextUser(ctx context.Context, user ContextUser) context.Context {
	return context.WithValue(ctx, contextUserKey, user)
}

/*
const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func SetUserIDIntoContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

*/
