package worker

import (
	"context"
	"sync"
)

type Job func(ctx context.Context) error

type Pool struct {
	sem chan struct{}
}

func New(concurrency int) *Pool {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &Pool{sem: make(chan struct{}, concurrency)}
}

func (p *Pool) Run(ctx context.Context, jobs []Job) []error {
	var wg sync.WaitGroup
	errs := make([]error, len(jobs))
	for i, job := range jobs {
		wg.Add(1)
		go func(idx int, fn Job) {
			defer wg.Done()
			p.sem <- struct{}{}
			defer func() { <-p.sem }()
			errs[idx] = fn(ctx)
		}(i, job)
	}
	wg.Wait()
	return errs
}
