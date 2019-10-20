# srv
Experimental quick setup of a http server in go with sane defaults.

[![GoDoc](https://godoc.org/github.com/karlpokus/srv?status.svg)](https://godoc.org/github.com/karlpokus/srv)

# features
- read, write and idle timeouts
- max header size
- bring your own router as long as it's a `http.Handler`
- graceful exit of server and any user supplied `Exiter`s

# usage
See cmd/main or `*test.go` for now

# tests
```bash
$ go test -v -race
```

# todos
- [x] use http.ServerMux as default router
- [ ] replace fmt with user supplied io.Writers
- [x] exiter interface on a graceful exit queue
- [ ] user supplied Exiters
- [ ] replace time lib with exponent notation

# license
MIT
