package srv

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/karlpokus/srv/testdata/routes"
)

func TestServeMux(t *testing.T) {
	conf := func(s *Server) error {
		router := http.NewServeMux()
		router.Handle("/hi", routes.Hello("bob"))
		s.Router = router
		s.StdoutWriter = ioutil.Discard
		return nil
	}
	_, err := New(conf)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}

func TestHttpRouter(t *testing.T) {
	conf := func(s *Server) error {
		router := httprouter.New()
		router.HandlerFunc("GET", "/greet/:user", routes.Greet)
		s.Router = router
		s.StdoutWriter = ioutil.Discard
		return nil
	}
	s, err := New(conf)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	r := httptest.NewRequest("GET", "/greet/bob", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, res.StatusCode)
	}
	if !bytes.Equal(bytes.TrimSpace(body), []byte("hello bob")) {
		t.Errorf("expected %s, got %s", []byte("hello bob"), bytes.TrimSpace(body))
	}
}
