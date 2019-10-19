package routes

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Hello(s string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello %s", s)
	}
}

func Greet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	fmt.Fprintf(w, "hello %s", params.ByName("user"))
}
