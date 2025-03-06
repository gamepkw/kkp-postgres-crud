package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"postgres-crud/app/internal/config"
	handler "postgres-crud/app/internal/handlers"
	repository "postgres-crud/app/internal/repositories"
	service "postgres-crud/app/internal/services"
	database "postgres-crud/app/pkg/postgres"

	validate "postgres-crud/app/pkg/validator"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()

	// Load config
	conf := config.Load(ctx)

	// Initialize database
	db, err := database.New(ctx, database.Config{
		Host:     conf.Secrets.PostgresHost,
		Port:     conf.Secrets.PostgresPort,
		User:     conf.Secrets.PostgresUser,
		Password: conf.Secrets.PostgresPassword,
		DBName:   conf.DB.Name,
	})
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot init database")
	}

	// Initialize repository, service, and handler
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	e := echo.New()

	v := validator.New()
	e.Validator = &validate.CustomValidator{Validator: v}

	// Define routes
	rootCtx := e.Group("/postgres-crud")
	api := rootCtx.Group("/api")
	v1 := api.Group("/v1")

	v1.POST("/orders", h.CreateOrder)
	v1.GET("/orders", h.GetAll)
	v1.GET("/orders/:order_id", h.GetByID)
	v1.PUT("/orders/:order_id/status", h.UpdateOrderStatus)

	// Start server
	go func() {
		if err := e.Start(conf.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
