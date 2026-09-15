package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

// backends are the upstream server instances to load-balance across.
// Reuses the servers from 02-scaling (run them with -port=8081/8082/8083).
var backends = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

// roundRobin picks the next backend in sequence on every call. next is
// accessed by every incoming request's goroutine concurrently, so it's
// incremented atomically rather than with a plain "counter++" (which
// would be a data race under concurrent requests).
type roundRobin struct {
	targets []*url.URL
	next    uint64
}

func newRoundRobin(rawURLs []string) *roundRobin {
	targets := make([]*url.URL, len(rawURLs))
	for i, raw := range rawURLs {
		u, err := url.Parse(raw)
		if err != nil {
			log.Fatalf("invalid backend url %q: %v", raw, err)
		}
		targets[i] = u
	}
	return &roundRobin{targets: targets}
}

func (rr *roundRobin) pick() *url.URL {
	n := atomic.AddUint64(&rr.next, 1)
	return rr.targets[(n-1)%uint64(len(rr.targets))]
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[proxy] %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	port := flag.String("port", "9090", "port for the proxy to listen on")
	flag.Parse()

	rr := newRoundRobin(backends)

	// Director rewrites each incoming request to point at the next
	// backend in rotation before ReverseProxy forwards it upstream.
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			target := rr.pick()
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			log.Printf("[proxy] routing %s -> %s", req.URL.Path, target.Host)
		},
	}

	addr := ":" + *port
	log.Printf("reverse proxy listening on %s, backends: %v", addr, backends)
	if err := http.ListenAndServe(addr, loggingMiddleware(proxy)); err != nil {
		log.Fatal(err)
	}
}
