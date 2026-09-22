package main

import (
	"fmt"
	"sync"
	"time"
)

type Result struct {
	Position int
	Square   int
}

func main() {
	start := time.Now()
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var wg sync.WaitGroup
	wg.Add(len(nums))

	resCh := make(chan Result)
	results := make([]Result, len(nums))

	for id, n := range nums {
		go func(n int) {
			defer wg.Done()
			resCh <- Result{Position: id, Square: square(n)}
		}(n)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	total := 0
	for res := range resCh {
		total += res.Square
		results[res.Position] = res
	}

	for id, n := range nums {
		fmt.Printf("%d : %d\n", n, results[id].Square)
	}

	fmt.Printf("Total: %d\n", total)
	fmt.Printf("Elapsed: %v\n", time.Since(start))
}

func square(n int) int {
	return n * n
}
