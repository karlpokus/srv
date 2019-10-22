# srv
Experimental quick setup of a http server in go with sane defaults.

[![GoDoc](https://godoc.org/github.com/karlpokus/srv?status.svg)](https://godoc.org/github.com/karlpokus/srv)

# features
- read, write and idle timeouts
- max header size
- bring your own router as long as it's a `http.Handler`. Default router is `http.ServeMux`
- graceful exit of server and any user supplied `Exiter`s

# usage
See cmd/main or `*test.go` for now

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
- [ ] use routest when v2 has been fixed

# license
MIT
