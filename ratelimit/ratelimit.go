package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Service interface {
	NewLimiter(frequency float64, burst int) Limiter
	Close()
}

func NewService(cleanupInterval, idleTimeout time.Duration) Service {
	s := &service{}
	s.wg.Add(1)
	go s.cleanup(cleanupInterval, idleTimeout)
	return &service{
		limiters: make([]*limiter, 0),
		done:     make(chan struct{}),
	}
}

type service struct {
	limiters []*limiter
	done     chan struct{}
	wg       sync.WaitGroup
	mutex    sync.Mutex
}

func (s *service) NewLimiter(frequency float64, burst int) Limiter {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	l := &limiter{
		entries:   make(map[string]*entry),
		frequency: frequency,
		burst:     burst,
	}
	s.limiters = append(s.limiters, l)
	return l
}

func (s *service) Close() {
	close(s.done)
	s.wg.Wait()
}

func (s *service) cleanup(interval, idleTimeout time.Duration) {
	for {
		select {
		case <-s.done:
			s.wg.Done()
		default:
			time.Sleep(interval)
			s.mutex.Lock()
			for _, l := range s.limiters {
				l.mutex.Lock()
				for key, e := range l.entries {
					if time.Since(e.touchedAt) > idleTimeout {
						delete(l.entries, key)
					}
				}
				l.mutex.Unlock()
			}
			s.mutex.Unlock()
		}

	}
}

type Limiter interface {
	Allow(string) bool
}

type limiter struct {
	entries   map[string]*entry
	frequency float64
	burst     int
	mutex     sync.Mutex
}

type entry struct {
	limiter   *rate.Limiter
	touchedAt time.Time
}

func (l *limiter) Allow(key string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	e, ok := l.entries[key]
	if !ok {
		limiter := rate.NewLimiter(rate.Limit(l.frequency), l.burst)
		l.entries[key] = &entry{limiter, time.Now()}
		return limiter.Allow()
	}
	e.touchedAt = time.Now()
	return e.limiter.Allow()
}
