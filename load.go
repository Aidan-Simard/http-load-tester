package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

type Stat struct {
	reqTime float64
	success bool
}

func main() {
	// setup flags
	url := flag.String("u", "", "URL to be tested")
	numReq := flag.Int("n", 1, "Number of requests to make")
	concReq := flag.Int("c", 1, "Number of concurrent requests")
	flag.Parse()

	if *url == "" {
		log.Fatalf("error: required url argument (-u) missing.")
	}

	// init channels and waitgroup for sync
	results := make(chan Stat, *numReq)
	var wg sync.WaitGroup
	wg.Add(2)

	// make requests and fetch results from channel
	go makeRequests(*numReq, *concReq, url, results, &wg)
	go getResults(results, &wg)

	// wait for all requests and results to finish
	wg.Wait()
}

func makeRequests(numReq int, concReq int, url *string, results chan Stat, wg *sync.WaitGroup) {
	var swg sync.WaitGroup
	swg.Add(concReq)

	for i := 0; i < concReq; i++ {
		go func(n int) {
			for j := 0; j < n; j++ {
				makeRequest(url, results)
			}
			swg.Done()
		}(numReq / concReq)
	}

	swg.Wait()

	close(results)
	wg.Done()
}

func makeRequest(url *string, results chan Stat) {
	start := time.Now()
	resp, err := http.Get(*url)
	end := time.Since(start)

	if err != nil {
		log.Fatalf("error: %s", err)
	}

	stat := Stat{reqTime: float64(end.Microseconds()), success: resp.StatusCode == http.StatusOK}

	results <- stat
}

func getResults(results chan Stat, wg *sync.WaitGroup) {
	successes := 0
	failures := 0
	minReqTime := math.MaxFloat64
	meanReqTime := 0.0
	maxReqTime := 0.0

	start := time.Now()

	for stat := range results {
		if stat.success {
			successes += 1
		} else {
			failures += 1
		}

		if stat.reqTime < minReqTime {
			minReqTime = stat.reqTime
		}

		if stat.reqTime > maxReqTime {
			maxReqTime = stat.reqTime
		}

		meanReqTime += stat.reqTime
	}

	meanReqTime = meanReqTime / (float64(successes + failures))
	total := time.Since(start)
	rps := float64(successes+failures) / total.Seconds()

	fmt.Printf("\nSuccesses: %d\nFailures: %d\nTotal Time: %d (ms)\nRequests Per Second: %d\n\nMin Request Time: %f (us)\nMean Request Time: %f (us)\nMax Request Time: %f (us)\n", successes, failures, total.Microseconds(), int(rps), minReqTime, meanReqTime, maxReqTime)
	wg.Done()
}
