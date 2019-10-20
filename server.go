// package srv provides a quick setup of a http server with sane defaults
package srv

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultHost       = "127.0.0.1"
	defaultPort       = "9012"
	gracePeriod int64 = 5e9
)

type Conf struct {
	Host, Port string
	ExiterList []Exiter
}

type Server struct {
	*http.Server
	Router http.Handler
	Conf
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

type ConfFunc func(*Server) error

// New passes a Server type to a user-supplied ConfFunc and returns
// a ready to use http server
func New(fn ConfFunc) (*Server, error) {
	s := &Server{}
	err := fn(s)
	if err != nil {
		return nil, err
	}
	s.Server = &http.Server{
		Addr:              addrWithDefaults(s.Host, s.Port),
		Handler:           s.Router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
	return s, nil
}

// Start runs a signal listener and starts the Server
func (s *Server) Start() error {
	fmt.Println("srv starting up")
	errc := make(chan error)
	go func() {
		// return this err to caller at startup and runtime
		// once we get an interrupt we ignore it
		errc <-s.ListenAndServe()
	}()
	ttl := time.Duration(1e9)
	// consider server start successful after ttl
	select {
	case err := <-errc:
		return fmt.Errorf("Error starting server: %s", err)
	case <-time.After(ttl):
		fmt.Printf("srv listening on %s\n", s.Server.Addr)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// select between server runtime err vs interrupt signal
	select {
	case err := <-errc:
		return fmt.Errorf("Error running server: %s", err)
	case <-sigc:
		return gracefulExit(gracePeriod, append(s.ExiterList, s.Server))
	}
}

// addrWithDefaults returns the default addr if vars are not set
func addrWithDefaults(host, port string) string {
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	return fmt.Sprintf("%s:%s", host, port)
}
