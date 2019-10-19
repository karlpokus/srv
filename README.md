# srv
Experimental quick setup of a http server in go with sane defaults.

[![GoDoc](https://godoc.org/github.com/karlpokus/srv?status.svg)](https://godoc.org/github.com/karlpokus/srv)

# features
- read and write timeouts
- max header size
- use whatever router you want
- graceful exit

# usage
See cmd/main

# tests
```bash
$ go test -v -race
```

# todos
- [x] use http.ServerMux as default router
- [ ] replace fmt with user supplied io.Writers
- [ ] user supplied "Exiter interface" to be placed on a graceful exit queue

# license
MIT
