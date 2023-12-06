package otlp

import (
	"context"

	"github.com/elastic/apm-data/input"
	"go.elastic.co/apm/v2"
)

func semAcquire(ctx context.Context, sem input.Semaphore, i int64) error {
	sp, ctx := apm.StartSpan(ctx, "Semaphore.Acquire", "Reporter")
	defer sp.End()

	return sem.Acquire(ctx, i)
}
