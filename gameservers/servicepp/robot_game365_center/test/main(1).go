package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/valyala/fasthttp"
)

const (
	url         = "https://game.ppply.site/Cash/Get"
	concurrency = 5000   // 并发数
	requests    = 100000 // 总请求数
)

func main() {
	fmt.Println("Testing fasthttp...")
	testFasthttp()

	fmt.Println("Testing go-resty...")
	testResty()
}

func testFasthttp() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			client := &fasthttp.Client{
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second}
			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			req.SetRequestURI(url)
			for j := 0; j < requests/concurrency; j++ {
				if err := client.Do(req, resp); err != nil || string(resp.Body()) != `{"code":1,"error":"Invalid App Params","data":null}` {
					fmt.Println("Error:", err)
					return
				}

			}
		}()
	}

	wg.Wait()
	fmt.Printf("fasthttp completed in %v\n", time.Since(start))
}

func testResty() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			client := resty.New()
			for j := 0; j < requests/concurrency; j++ {
				resp, err := client.R().Get(url)
				if err != nil || resp.String() != `{"code":1,"error":"Invalid App Params","data":null}` {
					fmt.Println("Error:", err)
					return
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("go-resty completed in %v\n", time.Since(start))
}
