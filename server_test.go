package srv

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/karlpokus/srv/testdata/routes"
	"github.com/karlpokus/routest/v2"
)

func TestDefaultRouter(t *testing.T) {
	routest.Test(t, func() http.Handler {
		s, err := New(func(s *Server) error {
			router := s.DefaultRouter()
			router.Handle("/hi", routes.Hello("bob"))
			s.Router = router
			s.Quiet()
			return nil
		})
		if err != nil {
			t.Errorf("expected nil, got %s", err)
		}
		return s
	}, []routest.Data{
		{
			"hello defaultrouter",
			"GET",
			"/hi",
			nil,
			nil,
			200,
			[]byte("hello bob"),
		},
	})
}

func TestCustomRouter(t *testing.T) {
	routest.Test(t, func() http.Handler {
		s, err := New(func(s *Server) error {
			router := httprouter.New()
			router.HandlerFunc("GET", "/greet/:user", routes.Greet)
			s.Router = router
			s.Quiet()
			return nil
		})
		if err != nil {
			t.Errorf("expected nil, got %s", err)
		}
		return s
	}, []routest.Data{
		{
			"greet customrouter",
			"GET",
			"/greet/lucy",
			nil,
			nil,
			200,
			[]byte("hello lucy"),
		},
	})
}
