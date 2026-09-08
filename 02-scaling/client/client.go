package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var ports = []string{"8081", "8082", "8083"}

// requestsPerInstance controls how many concurrent requests each server
// instance receives, so we can see how load spreads across horizontal
// copies vs a single vertically-scaled instance.
const requestsPerInstance = 5

func hit(port string, n int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	url := fmt.Sprintf("http://localhost:%s/hello", port)

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		results <- fmt.Sprintf("[:%s #%d] error: %v", port, n, err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	results <- fmt.Sprintf("[:%s #%d] %s in %s", port, n, resp.Status, time.Since(start))
}

func main() {
	var wg sync.WaitGroup
	results := make(chan string, len(ports)*requestsPerInstance)

	overallStart := time.Now()
	for _, port := range ports {
		for i := 1; i <= requestsPerInstance; i++ {
			wg.Add(1)
			go hit(port, i, &wg, results)
		}
	}

	wg.Wait()
	close(results)

	for r := range results {
		log.Println(r)
	}
	log.Printf("total wall time for %d concurrent requests across %d instances: %s",
		len(ports)*requestsPerInstance, len(ports), time.Since(overallStart))
}
