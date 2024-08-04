package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"Body"`
}

func endPointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")

	writer.WriteHeader(http.StatusOK)

	message := Message{
		Status: "Successful",
		Body:   "You have reached the API",
	}

	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return
	}
}

func main() {

	http.Handle("/ping", perUserRateLimiter(endPointHandler))

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Println("Something went wrong", err)
	}

}

func perUserRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(2, 4),
			}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {

			message := Message{
				Status: "FAILED",
				Body:   "limit exceeded",
			}
			writer.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(writer).Encode(&message)
			if err != nil {
				return
			}
			return
		}
		fmt.Printf("ip: %v\n", ip)
		fmt.Printf("clients[ip]: %v\n", clients[ip])
		mu.Unlock()
		next(writer, request)
	})

}
