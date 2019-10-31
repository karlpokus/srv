package main

import (
	"log"
	"os"

	"github.com/karlpokus/srv"
	"github.com/karlpokus/srv/testdata/routes"
)

var stdout = log.New(os.Stdout, "test ", 0)

func conf(s *srv.Server) error {
	router := s.DefaultRouter()
	router.Handle("/hi", routes.Hello("bob"))
	s.Router = router
	s.Logger = stdout
	return nil
}

func main() {
	s, err := srv.New(conf)
	if err != nil {
		stdout.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		stdout.Fatal(err)
	}
	stdout.Println("main exited")
}
