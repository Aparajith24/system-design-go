# system-design-go

Learning system design concepts one at a time, each implemented as a small,
runnable Go program. One topic per numbered folder, built from scratch using
just the Go standard library (`net/http`, `net/http/httputil`, `sync`, etc.)
— no frameworks.

## Topics

### 01 — Minimal HTTP server + client
`01-http-server-client/`

A basic `net/http` server and client. The server logs method, path, and
response time for every request via middleware; the client times its own
round trip.

```
go run ./01-http-server-client/server
go run ./01-http-server-client/client
```

### 02 — Vertical vs horizontal scaling
`02-scaling/`

Same server, now port-configurable (`-port` flag), so multiple identical
copies can run side by side (horizontal scaling). A client fires concurrent
requests at 3 instances using goroutines + `sync.WaitGroup`, timing each one
and the total wall-clock time, to see throughput scale with instance count.

```
go run ./02-scaling/server -port=8081
go run ./02-scaling/server -port=8082
go run ./02-scaling/server -port=8083
go run ./02-scaling/client
```

### 03 — Goroutines and real parallelism
`03-goroutines/`

A CPU-bound workload (naive prime counting over a range) run once
sequentially and once split across `NumCPU` goroutines with results merged
via `atomic.AddInt64`. Demonstrates actual multi-core parallelism, not just
non-blocking I/O concurrency.

```
go run ./03-goroutines
```

```

## Structure

Each topic lives in its own numbered directory and is a self-contained
`package main` (or `server`/`client`/`proxy` subpackages where a topic needs
multiple cooperating programs). Nothing is shared between topics beyond
reusing a previous topic's server as a building block, as noted above.
