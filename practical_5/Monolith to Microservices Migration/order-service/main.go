package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/database"
	"order-service/handlers"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	consulapi "github.com/hashicorp/consul/api"
)

func registerWithConsul(serviceName string, port int) error {
	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"

	consul, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()

	registration := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", serviceName, hostname),
		Name:    serviceName,
		Port:    port,
		Address: hostname,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", hostname, port),
			Interval: "10s",
			Timeout:  "3s",
		},
	}

	return consul.Agent().ServiceRegister(registration)
}

func main() {
	// Connect to dedicated order database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=order_db port=5432 sslmode=disable"
	}

	if err := database.Connect(dsn); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Add health endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Order endpoints (note: no /api prefix)
	r.Post("/orders", handlers.CreateOrder)
	r.Get("/orders", handlers.GetOrders)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Register with Consul
	portInt, _ := strconv.Atoi(port)
	if err := registerWithConsul("order-service", portInt); err != nil {
		log.Printf("Failed to register with Consul: %v", err)
	}

	log.Printf("Order service starting on :%s", port)
	http.ListenAndServe(":"+port, r)
}
