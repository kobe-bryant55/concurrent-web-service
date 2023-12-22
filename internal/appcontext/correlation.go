package appcontext

import (
	"context"
	"errors"

	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
)

var ErrInvalidRole = errors.New("unable to extract role from appcontext")

type mtsBlogRoleCtxKey struct{}

func WithRole(ctx context.Context, role types.Role) context.Context {
	return context.WithValue(ctx, mtsBlogRoleCtxKey{}, role)
}

func Role(ctx context.Context) (types.Role, error) {
	role, ok := ctx.Value(mtsBlogRoleCtxKey{}).(types.Role)
	if !ok {
		return "", ErrInvalidRole
	}

	return role, nil
}
