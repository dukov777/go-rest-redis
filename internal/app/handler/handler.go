package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type Payload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func PayloadHandler(client *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the raw body for debugging
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		bodyString := string(bodyBytes)
		log.Printf("Received payload: %s", bodyString)

		// Use the bytes to decode the payload
		var p Payload
		err = json.Unmarshal(bodyBytes, &p)
		if err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use the Redis client to store the data
		log.Printf("Storing key: %s, value: %s", p.Key, p.Value)
		err = client.Set(r.Context(), p.Key, p.Value, 0).Err()
		if err != nil {
			log.Printf("Error storing in Redis: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All key-value pairs received and stored in Redis.\n"))
	}
}
