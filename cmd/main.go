package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	"simple-order-service/internal/http"
	"simple-order-service/internal/repo"
	"simple-order-service/internal/service"
	"simple-order-service/internal/service/kafka"
	"simple-order-service/internal/service/redis"
	"simple-order-service/pkg/db"
	"time"
)

var pgRepo repo.OrderRepo
var redisClient *redis.Client
var orderService *service.OrderService
var router netHttp.Handler
var orderHandler *http.OrderHandler

func main() {

	pgDb, err := db.ConnectToDb()
	if err != nil {
		log.Fatal(err)
	}

	pgRepo = repo.NewPgOrderRepo(pgDb)

	redisClient = redis.NewRedisClient("127.0.0.1:6379", 5*time.Minute, pgRepo)

	orderService = service.NewOrderService(pgRepo, redisClient)

	if err := redisClient.RestoreCache(context.Background()); err != nil {
		log.Fatal(err)
	}

	StartKafka()

	orderHandler = http.NewOrderHandler(*orderService)
	router = http.NewRouter(orderHandler)
	StartWebServer()
}

func StartWebServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting HTTP server on port %s", port)
	if err := netHttp.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func StartKafka() {
	brokers := []string{"localhost:9092"}

	groupID := "order-service-group"
	topic := "orders"

	go func() {
		consumer := kafka.NewConsumer(pgRepo, redisClient, brokers, groupID, topic)
		consumer.Consume()
	}()
}
