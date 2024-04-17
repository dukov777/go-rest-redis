package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	redis "my-go-project/pkg/redis"
	"net/http"
)

type Payload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Update the function signature to accept the gRPC client
func PayloadHandler(redisClient redis.RedisClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		var p Payload
		if err := json.Unmarshal(bodyBytes, &p); err != nil {
			log.Printf("Error decoding JSON: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Use the Redis client to store the data
		log.Printf("Storing key: %s, value: %s", p.Key, p.Value)
		if err := redisClient.SetKey(r.Context(), p.Key, p.Value, 0); err != nil {
			log.Printf("Error storing in Redis: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("All key-value pairs received and stored in Redis.\n"))
	}
}
