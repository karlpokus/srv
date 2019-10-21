package srv

import (
	"context"
	"errors"
	"time"
)

type Exiter interface {
	Shutdown(context.Context) error
}

var (
	ExitErr     = errors.New("Graceful shutdown completed with errors")
	ExitTimeout = errors.New("Graceful shutdown timeout")
)

// gracefulExit shuts down Exiters gracefully and returns an error if any
func gracefulExit(graceperiod int64, queue []Exiter) error {
	stdout.Println("Graceful shutdown start")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(graceperiod))
	defer cancel()

	errc := make(chan error)
	resc := make(chan []error)
	go func() {
		out := make([]error, len(queue), len(queue))
		for i := 0; i < len(queue); i++ {
			out[i] = <-errc
		}
		resc <- out
	}()
	for _, q := range queue {
		go func(q Exiter) {
			errc <-q.Shutdown(ctx)
		}(q)
	}
	select {
	case res := <-resc:
		if hasErrs(res) {
			return ExitErr
		}
		stdout.Println("Graceful shutdown complete")
		return nil
	case <-ctx.Done(): // this chan closed recieve will not block other recieves
		return ExitTimeout
	}
}

// hasErrs determines if the list contains at least one error
func hasErrs(list []error) bool {
	for _, err := range list {
		if err != nil {
			return true
		}
	}
	return false
}
