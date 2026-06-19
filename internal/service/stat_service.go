package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/FoPQer/go-shortener/internal/model"
	"golang.org/x/sync/errgroup"
)

type urlsCounter interface {
	Count(ctx context.Context) (int, error)
}

type usersCounter interface {
	Count(ctx context.Context) (int, error)
}

type StatService struct {
	urlsCounter  urlsCounter
	usersCounter usersCounter
}

func NewStatService(urlsCounter urlsCounter, usersCounter usersCounter) *StatService {
	return &StatService{
		urlsCounter:  urlsCounter,
		usersCounter: usersCounter,
	}
}

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
