package main

import (
	"context"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kimvlry/simple-order-service/internal/config"
	"github.com/kimvlry/simple-order-service/internal/http"
	"github.com/kimvlry/simple-order-service/internal/interfaces"
	"github.com/kimvlry/simple-order-service/internal/repo"
	"github.com/kimvlry/simple-order-service/internal/service"
	"github.com/kimvlry/simple-order-service/internal/service/kafka"
	"github.com/kimvlry/simple-order-service/internal/service/redis"
	"github.com/kimvlry/simple-order-service/pkg/db"
)

type App struct {
	config       *config.Config
	pgRepo       interfaces.OrderRepository
	redisClient  interfaces.Cache
	orderService interfaces.OrderService
	consumer     interfaces.MessageConsumer
	server       *netHttp.Server
}

func main() {
	cfg := config.Load()

	app := &App{
		config: cfg,
	}

	if err := app.initialize(); err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	app.start()

	// Graceful shutdown
	app.gracefulShutdown()
}

func (a *App) initialize() error {
	// Initialize database
	pgDb, err := db.ConnectToDb(a.config.Database)
	if err != nil {
		return err
	}

	// Initialize repositories and services
	a.pgRepo = repo.NewPgOrderRepo(pgDb)
	a.redisClient = redis.NewRedisClient(a.config.Redis, a.pgRepo)
	a.orderService = service.NewOrderService(a.pgRepo, a.redisClient)

	// Restore cache
	if err := a.redisClient.RestoreCache(context.Background()); err != nil {
		return err
	}

	// Initialize Kafka consumer
	a.consumer = kafka.NewConsumer(a.pgRepo, a.redisClient, a.config.Kafka)

	// Initialize HTTP server
	orderHandler := http.NewOrderHandler(a.orderService)
	router := http.NewRouter(orderHandler)
	a.server = &netHttp.Server{
		Addr:    ":" + a.config.Server.Port,
		Handler: router,
	}

	return nil
}

func (a *App) start() {
	// Start Kafka consumer
	go func() {
		if err := a.consumer.Consume(context.Background()); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Start HTTP server
	go func() {
		log.Printf("Starting HTTP server on port %s", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != netHttp.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func (a *App) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close consumer
	if err := a.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	// Close cache
	if err := a.redisClient.Close(); err != nil {
		log.Printf("Error closing cache: %v", err)
	}

	// Close database
	if err := a.pgRepo.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited")
}
