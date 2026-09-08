package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(port string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[:%s] %s %s %s", port, r.Method, r.URL.Path, time.Since(start))
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// simulate a bit of work so concurrent timings are visible
	time.Sleep(100 * time.Millisecond)
	w.Write([]byte("hello, world\n"))
}

func main() {
	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	addr := ":" + *port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(*port, mux)); err != nil {
		log.Fatal(err)
	}
}
