package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	requests := make(chan int, 100)
	for i := 1; i <= 100; i++ {
		requests <- i
	}
	close(requests)
	limiter := time.Tick(200 * time.Millisecond)

	for req := range requests {
		<-limiter
		res, err := http.Get("http://localhost:8080/ratelimit")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
		fmt.Println("request", req, time.Now())
	}
}
