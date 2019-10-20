package main

import (
	"net/http"
	"log"

	"github.com/karlpokus/srv"
	"github.com/karlpokus/srv/testdata/routes"
)

func conf(s *srv.Server) error {
	router := http.NewServeMux()
	router.Handle("/hi", routes.Hello("bob"))
	s.Router = router
	return nil
}

func main() {
	s, err := srv.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("main exited")
}
