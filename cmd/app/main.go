package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mijgona/salon-crm/cmd"
	"github.com/mijgona/salon-crm/internal/adapters/out/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx := context.Background()

	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://salon_crm:salon_crm_pass@localhost:5432/salon_crm_db?sslmode=disable"
	}

	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL")

	// Initialize composition root with pgx pool
	cr := cmd.NewCompositionRoot(pool)

	// Initialize Mediatr (registers event handlers)
	_ = cr.Mediatr()

	// Create Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Register HTTP handlers
	cr.ClientHTTPHandler().Register(e)
	cr.AppointmentHTTPHandler().Register(e)
	cr.LoyaltyHTTPHandler().Register(e)
	cr.CalendarHTTPHandler().Register(e)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting salon-crm server on port %s", port)
	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
