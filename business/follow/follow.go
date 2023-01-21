package follow

import (
	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
	"github.com/samuelsih/guwu/model"
	"github.com/samuelsih/guwu/pkg/redis"
	"github.com/samuelsih/guwu/pkg/securer"
)

type Deps struct {
	DB *sqlx.DB
	GetUserSession func(ctx context.Context, key string, dst any) error
}

type FollowIn struct {
	UserID string `json:"user_id"`
}

type FollowOut struct {
	business.CommonResponse
}

func(d *Deps) Follow(ctx context.Context, in FollowIn, common business.CommonInput) FollowOut {
	var out FollowOut
	user := model.User{}

	sessionID, err := securer.Decrypt(common.SessionID)
	if err != nil {
		out.SetError(403, "Unauthenticated")
		return out
	}

	err = d.GetUserSession(ctx, string(sessionID), &user)
	if err != nil {
		if errors.Is(err, redis.ErrInternal) {
			out.SetError(500, redis.ErrInternal.Error())
			return out
		}

		out.SetError(400, err.Error())
		return out
	}

	log.Println("ini user:", user)

	userFollow := model.NewUserFollow(d.DB)

	if !userFollow.Insert(ctx, user.ID, in.UserID) {
		out.SetError(userFollow.StatusCode, userFollow.Err.Error())
		return out
	}

	out.SetOK()
	return out
}

func (d *Deps) Unfollow(ctx context.Context) {

}