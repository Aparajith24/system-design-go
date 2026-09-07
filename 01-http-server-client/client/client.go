package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080/hello"

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GET %s -> %s in %s", url, resp.Status, time.Since(start))
	log.Printf("body: %s", body)
}
