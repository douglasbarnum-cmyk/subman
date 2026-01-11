package main

import (
	"log"

	"subman/internal/service"
	"subman/internal/storage"
	"subman/internal/ui"
)

func main() {
	// Initialize storage
	store, err := storage.NewJSONStorage()
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize services
	svc := service.NewSubscriptionService(store)
	paymentSvc := service.NewPaymentService(store)

	// Create and run UI
	app := ui.NewApp(svc, paymentSvc)
	app.Run()
}
