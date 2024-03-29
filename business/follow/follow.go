package follow

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/model"

	"github.com/samuelsih/guwu/pkg/securer"
)

type Deps struct {
	DB             *sqlx.DB
	GetUserSession func(ctx context.Context, key string, dst any) error
}

type FollowIn struct {
	UserID string `json:"user_id"`
}

type FollowOut struct {
	business.CommonResponse
}

func (d *Deps) Follow(ctx context.Context, in FollowIn, common business.CommonInput) FollowOut {
	var out FollowOut
	var user model.User

	if common.SessionID == "" {
		out.RawError(403, "Unauthenticated")
		return out
	}

	sessionID, err := securer.Decrypt(common.SessionID)
	if err != nil {
		out.SetError(err)
		return out
	}

	err = d.GetUserSession(ctx, string(sessionID), &user)
	if err != nil {
		out.SetError(err)
		return out
	}

	err = model.FollowUser(ctx, d.DB, user.ID, in.UserID)
	if err != nil {
		out.SetError(err)
		return out
	}

	out.SetOK()
	return out
}

type UnfollowIn struct {
	UserID string `json:"user_id"`
}

type UnfollowOut struct {
	business.CommonResponse
}

func (d *Deps) Unfollow(ctx context.Context, in FollowIn, common business.CommonInput) UnfollowOut {
	var out UnfollowOut
	var user model.User

	if common.SessionID == "" {
		out.RawError(403, "Unauthenticated")
		return out
	}

	sessionID, err := securer.Decrypt(common.SessionID)
	if err != nil {
		out.SetError(err)
		return out
	}

	err = d.GetUserSession(ctx, string(sessionID), &user)
	if err != nil {
		out.SetError(err)
		return out
	}

	err = model.UnfollowUser(ctx, d.DB, user.ID, in.UserID)
	if err != nil {
		out.SetError(err)
		return out
	}

	out.SetOK()
	return out
}
