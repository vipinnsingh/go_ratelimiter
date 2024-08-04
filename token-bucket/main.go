package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Status string
	Body   string
}

func endPointHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-type", "application/json")

	rw.WriteHeader(http.StatusOK)

	response := response{
		Status: "SUCCESS",
		Body:   "Hi You have successfully hit the api",
	}
	err := json.NewEncoder(rw).Encode(&response)
	if err != nil {
		return
	}
}

func main() {

	http.HandleFunc("/ping", rateLimiter(endPointHandler))

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Println("there was an error listening to the port 3000", err)
	}

}
