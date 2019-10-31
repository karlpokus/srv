package srv

import (
	"context"
	"fmt"
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
	<-ctx.Done()
	return nil
}

func TestGracefulExit(t *testing.T) {
	s, _ := New(func(s *Server) error {
		s.ExiterList = append(s.ExiterList, exiterMock{})
		return nil
	})
	err := gracefulExit(s.GracePeriod, append(s.ExiterList, s.Server))
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}

func TestGracefulExitErr(t *testing.T) {
	s, _ := New(func(s *Server) error {
		s.ExiterList = append(s.ExiterList, exiterMockErr{})
		return nil
	})
	err := gracefulExit(s.GracePeriod, append(s.ExiterList, s.Server))
	if err != ExitErr {
		t.Errorf("expected %s, got %s", ExitErr, err)
	}
}

func TestGracefulExitTimeout(t *testing.T) {
	s, _ := New(func(s *Server) error {
		s.ExiterList = append(s.ExiterList, exiterMockTimeout{}, exiterMockTimeout{})
		s.GracePeriod = "0"
		return nil
	})
	err := gracefulExit(s.GracePeriod, append(s.ExiterList, s.Server))
	if err != ExitTimeout {
		t.Errorf("expected %s, got %s", ExitTimeout, err)
	}
}
