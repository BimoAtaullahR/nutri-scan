package user

import "context"

func WithTestAnonymousUser(ctx context.Context, u AnonymousUser) context.Context {
	return context.WithValue(ctx, anonymousUserContextKey{}, u)
}
