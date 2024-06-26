package main

import (
	"log"
	"my-go-project/internal/app/handler"
	"my-go-project/internal/app/middleware"
	"my-go-project/pkg/redis"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting server...")
	redisClient := redis.NewClient()

	if redisClient == nil {
		log.Fatalf("Error connecting to Redis")
	}
	log.Default().Printf("Connected to redis")

	r := mux.NewRouter()
	// Use both middleware
	r.Use(middleware.TracingMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// Apply the logging middleware to all routes
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/", handler.PayloadHandler(redisClient)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
