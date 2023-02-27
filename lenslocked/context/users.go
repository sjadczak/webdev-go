package context

import (
	"context"

	"github.com/sjadczak/webdev-go/lenslocked/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	v := ctx.Value(userKey)
	user, ok := v.(*models.User)
	if !ok {
		return nil
	}
	return user
}
