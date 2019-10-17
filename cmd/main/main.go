package main

import (
	"github.com/karlpokus/srv"
	"github.com/karlpokus/srv/testdata/routes"
)

func conf(s *srv.Server) error {
	s.Router.Handle("/hi", routes.Hello("bob"))
	return nil
}

func main() {
	s, err := srv.New(conf)
	if err != nil {
		panic(err)
	}
	s.Start()
}
