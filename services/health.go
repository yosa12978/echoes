package services

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/yosa12978/echoes/data"
	"github.com/yosa12978/echoes/logging"
)

type HealthService interface {
	Healthcheck(ctx context.Context) error
}

type healthService struct {
	logger logging.Logger
	pings  []data.Pinger
}

func NewHealthService(logger logging.Logger, pings ...data.Pinger) HealthService {
	return &healthService{
		logger: logger,
		pings:  pings,
	}
}

func (s *healthService) Healthcheck(ctx context.Context) error {
	numOfChecks := len(s.pings)
	errCh := make(chan error, numOfChecks)
	var wg sync.WaitGroup
	for i := 0; i < numOfChecks; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if err := s.pings[n].Ping(ctx); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	if len(errCh) == 0 {
		return nil
	}

	errsStr := []string{}
	for err := range errCh {
		if err.Error() != "" {
			s.logger.Error(err.Error())
			errsStr = append(errsStr, err.Error())
		}
	}
	mergedErrorStr := strings.Join(errsStr, "\r\n\r\n")
	return errors.New(mergedErrorStr)
}
