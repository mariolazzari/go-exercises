package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Position int
	URL      string
	Status   int
	Elapsed  time.Duration
	Err      error
}

func main() {
	start := time.Now()

	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://go.dev",
		"https://does-not-exist.example",
	}
	size := len(urls)

	results := make([]Result, size)

	var wg sync.WaitGroup
	wg.Add(size)

	resCh := make(chan Result, size)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for i, url := range urls {
		go func() {
			start := time.Now()
			defer wg.Done()

			status, err := checkURL(url, &client)

			resCh <- Result{
				Position: i,
				URL:      url,
				Status:   status,
				Elapsed:  time.Since(start),
				Err:      err,
			}
		}()
	}

	// cleanup
	go func() {
		wg.Wait()
		close(resCh)
	}()

	for res := range resCh {
		results[res.Position] = res
	}

	for _, res := range results {
		fmt.Println(res)
	}

	fmt.Printf("Elapsed: %v\n", time.Since(start))
}

func checkURL(url string, client *http.Client) (int, error) {
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
