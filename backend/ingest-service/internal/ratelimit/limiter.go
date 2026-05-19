package ratelimit

import (
	"context"

	"golang.org/x/time/rate"
)

type Limiter struct {
	limiters map[string]*rate.Limiter
}

func New(perSecond float64, burst int) *Limiter {
	return &Limiter{
		limiters: map[string]*rate.Limiter{
			"facebook": rate.NewLimiter(rate.Limit(perSecond), burst),
			"youtube":  rate.NewLimiter(rate.Limit(perSecond), burst),
		},
	}
}

func (l *Limiter) Wait(ctx context.Context, provider string) error {
	lim, ok := l.limiters[provider]
	if !ok {
		return nil
	}
	return lim.Wait(ctx)
}
