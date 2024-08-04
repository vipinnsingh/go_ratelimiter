package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimiter(next func(rw http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !limiter.Allow() {
			response := response{
				Status: "FAILED",
				Body:   "limit exceeded",
			}
			rw.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(rw).Encode(&response)
			if err != nil {
				return
			}
			return
		} else {
			next(rw, req)
		}
	})

}
