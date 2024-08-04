package main

import (
	"encoding/json"
	"log"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
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

	message := Message{
		Status: "Request Failed",
		Body:   "You have reached the API req limit",
	}
	jsonMessage, _ := json.Marshal(message)

	tollboothLimiter := tollbooth.NewLimiter(1, nil)
	tollboothLimiter.SetMessageContentType("application/json")
	tollboothLimiter.SetMessage(string(jsonMessage))

	http.Handle("/ping", tollbooth.LimitFuncHandler(tollboothLimiter, endPointHandler))

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println("request failed with err")
	}
}
