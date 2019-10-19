package srv

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/karlpokus/srv/testdata/routes"
)

func conf(s *Server) error {
	// replace router
	router := httprouter.New()
	router.HandlerFunc("GET", "/greet/:user", routes.Greet)
	s.Router = router
	return nil
}

func TestServer(t *testing.T) {
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
