package auth

import (
	"context"

	"github.com/dmasior/service-go/internal/domain"
)

func NewContext(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

func FromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(userKey{}).(domain.User)
	return user, ok
}

func MustFromContext(ctx context.Context) domain.User {
	user, ok := FromContext(ctx)
	if !ok {
		panic("user not found in context")
	}
	return user
}
