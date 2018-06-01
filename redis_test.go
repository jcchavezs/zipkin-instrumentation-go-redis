package zipkinredis_test

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/alicebob/miniredis"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/recorder"
	"testing"
	"github.com/jcchavezs/zipkin-instrumentation-go-redis"
)

const key = "key"
const value = "value"

func TestSpansAreCreated(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	rep := recorder.NewReporter()

	tracer, err := zipkin.NewTracer(rep)

	opts := &redis.Options{
		Addr: s.Addr(),
	}
	c := redis.NewClient(opts)
	c.WrapProcess(zipkinredis.ProcessWithTracing(tracer))
	c.WithContext(context.Background())

	c.Set(key, value, 0)

	res := c.Get(key)
	if want, have := value, res.Val(); want != have {
		t.Errorf("unexpected output value, wanted %q, got %q", want, have)
	}

	spans := rep.Flush()

	if want, have := 2, len(spans); want != have {
		t.Errorf("unexpected number of spans, wanted %d, got %d", want, have)
	}
}