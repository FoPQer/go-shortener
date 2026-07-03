package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/FoPQer/go-shortener/internal/model"
	"golang.org/x/sync/errgroup"
)

// Counter describes minimal counting behavior required by StatService.
type Counter interface {
	Count(ctx context.Context) (int, error)
}

// StatService aggregates service-wide counters such as URLs and users.
type StatService struct {
	urlsCounter  Counter
	usersCounter Counter
}

// NewStatService creates a StatService that uses provided counters.
func NewStatService(urlsCounter Counter, usersCounter Counter) *StatService {
	return &StatService{
		urlsCounter:  urlsCounter,
		usersCounter: usersCounter,
	}
}

// GetStats concurrently fetches URL and user counts and returns aggregated statistics.
func (s *StatService) GetStats(ctx context.Context) (*model.Stat, error) {
	if s.urlsCounter == nil || s.usersCounter == nil {
		return nil, fmt.Errorf("stat service is not configured")
	}

	result := model.NewStat(0, 0)
	var mu sync.Mutex

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		count, err := s.urlsCounter.Count(gctx)
		if err != nil {
			return fmt.Errorf("failed to count urls: %w", err)
		}

		mu.Lock()
		result.IncrementURLs(count)
		mu.Unlock()

		return nil
	})

	g.Go(func() error {
		count, err := s.usersCounter.Count(gctx)
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}

		mu.Lock()
		result.IncrementUsers(count)
		mu.Unlock()

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
