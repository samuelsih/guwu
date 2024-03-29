package health

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samuelsih/guwu/business"
)

type Deps struct {
	DB *sqlx.DB
}

type HealthCheckOutput struct {
	business.CommonResponse
	PostgreStatus string `json:"postgre_status"`
	Time int64 `json:"time_now"`
}

func (d *Deps) Check(ctx context.Context, data business.CommonInput) HealthCheckOutput {
	var out HealthCheckOutput

	err := d.DB.PingContext(ctx)
	if err != nil {
		out.PostgreStatus = err.Error()
		return out
	}

	out.PostgreStatus = "OK"
	out.Time = time.Now().Unix()
	out.SetOK()

	return out
}
