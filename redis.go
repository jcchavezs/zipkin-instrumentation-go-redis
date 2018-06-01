package zipkinredis

import (
	"github.com/go-redis/redis"
	"context"
	"github.com/openzipkin/zipkin-go"
)

func ProcessWithTracing(tracer *zipkin.Tracer) func(oldProcess func(ctx context.Context, cmd redis.Cmder) error) func(ctx context.Context, cmd redis.Cmder) error  {
	return func(oldProcess func(context.Context, redis.Cmder) error) func(context.Context, redis.Cmder) error {
		return func(ctx context.Context, cmd redis.Cmder) error {
			span, _ := tracer.StartSpanFromContext(ctx, cmd.Name())
			defer span.Finish()

			err := oldProcess(ctx, cmd)
			if err != nil {
				zipkin.TagError.Set(span, err.Error())
			}

			return err
		}
	}
}