package sqllogger

import (
	"context"
	"time"
)

type Timing struct {
	Start time.Time
	End   time.Time
}

type timingKey struct{}

func WithTiming(ctx context.Context, timing Timing) context.Context {
	return context.WithValue(ctx, timingKey{}, timing)
}

func GetTiming(ctx context.Context) (Timing, bool) {
	timing := ctx.Value(timingKey{})
	if timing == nil {
		return Timing{}, false
	}
	return timing.(Timing), true
}
