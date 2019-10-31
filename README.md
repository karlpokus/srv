# srv
Experimental quick setup of a http server in go with sane defaults.

[![GoDoc](https://godoc.org/github.com/karlpokus/srv?status.svg)](https://godoc.org/github.com/karlpokus/srv)

# features
- read, write and idle timeouts
- max header size
- bring your own router as long as it's a `http.Handler`. Default router is `http.ServeMux`
- graceful exit of server and any user supplied `Exiter`s

# usage
Create a `srv.ConfFunc` type where you add your prefered server configuration and pass it to `srv.New`
```go
import (
	"log"
	"os"

	"github.com/karlpokus/srv"
)

func conf(s *srv.Server) error {
	// Add a router (type http.Handler - required)
	// We're using the default here which is a http.ServeMux.
	// See tests for how to add a custom router.
	router := s.DefaultRouter()
	router.Handle("/hi", someRoute)
	s.Router = router
	// Add logging (type *log.logger - optional)
	// srv only logs during upstart and before exit. Omit to keep srv quiet.
	s.Logger = log.New(os.Stdout, "server ", 0)
	// Set http host and port (type string - optional)
	// Omit to use default 127.0.0.1:9012
	s.Host = "0.0.0.0"
	s.Port = "12900"
	// Add exiters (type srv.Exiter - optional)
	// These will be run on interrupt signals to allow for graceful exits.
	// The server itself is included by default
	s.ExiterList = append(s.ExiterList, &newDB{})
	// Set gracePeriod (type string - optional)
	// Default is 5s
	s.GracePeriod = "0"
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

```

# tests
```bash
$ go test -v -race
```

# todos
- [x] replace fmt with user supplied io.Writers
- [x] exiter interface on a graceful exit queue
- [x] user supplied Exiters
- [ ] replace time lib with exponent notation
- [ ] consider testing Start func
- [ ] how to handle exit errors?
- [x] convenience funcs like Quiet, DefaultRouter
- [ ] log request opt
- [ ] optional endpoint to read ConnStates
- [x] expose Graceperiod
- [x] use routest when v2 has been fixed
- [x] conf should take a log.Logger instead of an io.Writer

# license
MIT
