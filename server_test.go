package srv

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/karlpokus/srv/testdata/routes"
)

func conf(s *Server) error {
	s.Router.Handle("/hi", routes.Hello("bob"))
	return nil
}

func TestServer(t *testing.T) {
	s, err := New(conf)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	r := httptest.NewRequest("GET", "/hi", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 200 {
		t.Errorf("expected %d, got %d", 200, res.StatusCode)
	}
	if !bytes.Equal(bytes.TrimSpace(body), []byte("bob")) {
		t.Errorf("expected %s, got %s", []byte("bob"), bytes.TrimSpace(body))
	}
}
