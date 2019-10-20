package srv

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Exiter interface {
	Shutdown(context.Context) error
}

var (
	ExitErr     = errors.New("Graceful shutdown complete with errors")
	ExitTimeout = errors.New("Graceful shutdown timeout")
)

// gracefulExit shuts down Exiters gracefully
func gracefulExit(graceperiod int64, queue []Exiter) error {
	fmt.Println("Graceful shutdown start")
	ttl := time.Duration(graceperiod)
	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()

	errc := make(chan error)
	resc := make(chan []error) // temp
	go func() {
		out := make([]error, 0, len(queue))
		for i := 0; i < len(queue); i++ {
			//out[i] <-errc
			out = append(out, <-errc)
		}
		resc <- out
	}()
	for _, q := range queue {
		go func(q Exiter) {
			errc <- q.Shutdown(ctx)
		}(q)
	}
	select {
	case res := <-resc:
		if hasErrs(res) {
			return ExitErr
		}
		fmt.Println("Graceful shutdown complete")
		return nil
	case <-time.After(ttl):
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
