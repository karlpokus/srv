# srv
Experimental quick setup of a http server in go with sane defaults.

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
- [ ] replace fmt with buffered log

# license
MIT
