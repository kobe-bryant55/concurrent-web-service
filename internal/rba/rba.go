package rba

import (
	"context"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/appcontext"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/types"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"
)

type RBA interface {
	CheckHasRole(ctx context.Context, requiredRole types.Role) error
}

var roleScores = map[types.Role]uint8{
	types.Admin:      15,
	types.Registered: 13,
}

func New() RBA {
	return &rba{}
}

type rba struct{}

func (r *rba) CheckHasRole(ctx context.Context, requiredRole types.Role) error {
	userRole, err := appcontext.Role(ctx)
	if err != nil {
		return errorutils.ErrUnauthorized
	}

	if roleScores[userRole] < roleScores[requiredRole] {
		return errorutils.ErrUnauthorized
	}

	return nil
}
