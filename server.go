// package srv provides a quick setup of a http server with sane defaults
package srv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Conf struct {
	Host, Port string
}

type Server struct {
	*http.Server
	Router *http.ServeMux
	Conf
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

type ConfFunc func(*Server) error

// New passes a Server type to a user-supplied ConfFunc and returns
// a ready to use http server
func New(fn ConfFunc) (*Server, error) {
	s := &Server{
		Router: http.NewServeMux(),
	}
	err := fn(s)
	if err != nil {
		return nil, err
	}
	s.Server = &http.Server{
		Addr:              addrWithDefaults(s.Host, s.Port),
		Handler:           s.Router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}
	return s, nil
}

// Start runs a signal listener and starts the Server
func (s *Server) Start() error {
	go gracefulExit(s)
	fmt.Printf("srv listening on %s\n", s.Server.Addr)
	return s.ListenAndServe()
}

// gracefulExit shuts down the server gracefully on SIGINT
func gracefulExit(s *Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	fmt.Println("Shutdown starting..")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown err: %s", err)
		return
	}
	fmt.Println("Shutdown complete")
}

// addrWithDefaults returns the default addr if vars are not set
func addrWithDefaults(host, port string) string {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "9012"
	}
	return fmt.Sprintf("%s:%s", host, port)
}
