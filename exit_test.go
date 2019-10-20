package srv

import (
	"context"
	"fmt"
	"time"
	"testing"
)

type exiterMock struct{}

func (m exiterMock) Shutdown(ctx context.Context) error {
	return nil
}

type exiterMockErr struct{}

func (m exiterMockErr) Shutdown(ctx context.Context) error {
	return fmt.Errorf("oops!")
}

type exiterMockTimeout struct{}

func (m exiterMockTimeout) Shutdown(ctx context.Context) error {
	time.Sleep(1e9)
	return nil
}

func TestGracefulExit(t *testing.T) {
	s, _ := New(func(s *Server) error {
		return nil
	})
	err := gracefulExit(gracePeriod, []Exiter{s.Server, exiterMock{}})
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}

func TestGracefulExitErr(t *testing.T) {
	s, _ := New(func(s *Server) error {
		return nil
	})
	err := gracefulExit(gracePeriod, []Exiter{s.Server, exiterMockErr{}})
	if err != ExitErr {
		t.Errorf("expected %s, got %s", ExitErr, err)
	}
}

func TestGracefulExitTimeout(t *testing.T) {
	s, _ := New(func(s *Server) error {
		return nil
	})
	var ttl int64 = 0
	err := gracefulExit(ttl, []Exiter{s.Server, exiterMockTimeout{}})
	if err != ExitTimeout {
		t.Errorf("expected %s, got %s", ExitTimeout, err)
	}
}
