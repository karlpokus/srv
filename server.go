// package srv provides a quick setup of a http server with sane defaults
package srv

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultHost        = "127.0.0.1"
	defaultPort        = "9012"
	defaultGracePeriod = "5s"
)

type Conf struct {
	Host, Port  string
	ExiterList  []Exiter
	GracePeriod string
	Logger      *log.Logger
}

type Server struct {
	*http.Server
	Router http.Handler
	Conf
}

// the Server type is also a http.Handler which makes for convenient testing
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// DefaultRouter returns a http.ServeMux
func (s *Server) DefaultRouter() *http.ServeMux {
	return http.NewServeMux()
}

type ConfFunc func(*Server) error

var logger = log.New(ioutil.Discard, "", 0)

// New runs the ConfFunc, set some defaults and returns a
// ready to use http server
func New(fn ConfFunc) (*Server, error) {
	s := &Server{}
	err := fn(s)
	if err != nil {
		return nil, err
	}
	if s.Logger != nil {
		logger = s.Logger
	}
	if s.GracePeriod == "" {
		s.GracePeriod = defaultGracePeriod
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
	logger.Println("Server created")
	return s, nil
}

// Start runs a signal listener and starts the Server
func (s *Server) Start() error {
	logger.Println("Server starting")
	errc := make(chan error)
	go func() {
		// return this err to caller at startup and runtime
		// once we get an interrupt we ignore it
		errc <- s.ListenAndServe()
	}()
	ttl := time.Duration(1e9)
	// consider server start successful after ttl
	select {
	case err := <-errc:
		return fmt.Errorf("Error starting server: %s", err)
	case <-time.After(ttl):
		logger.Printf("Server listening on %s\n", s.Server.Addr)
	}
	// select between server runtime err vs interrupt signal
	select {
	case err := <-errc:
		return fmt.Errorf("Error running server: %s", err)
	case <-interrupt():
		return gracefulExit(s.GracePeriod, append(s.ExiterList, s.Server))
	}
}

// interrupt returns a chan that recieves interrupt signals
func interrupt() <-chan os.Signal {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return sigc
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
