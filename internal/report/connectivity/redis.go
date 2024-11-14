package connectivity

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	types "github.com/accuknox/rinc/types/connectivity"

	"github.com/redis/go-redis/v9"
)

// Report reports the redis connectivity status.
func (r Reporter) redisReport(ctx context.Context) (*types.Redis, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	c := redis.NewClient(&redis.Options{
		Addr: r.conf.Redis.Addr,
	})
	pong, err := c.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("pinging server: %w", err)
	}
	slog.LogAttrs(
		ctx,
		slog.LevelDebug,
		"pinging redis server",
		slog.String("pong", pong),
	)
	return &types.Redis{
		Reachable: true,
	}, nil
}
