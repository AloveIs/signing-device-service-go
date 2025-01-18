package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

const (
	ListenAddress = ":8080"
)

func main() {
	// create the repository (database)
	repo := persistence.NewInMmemoryDb()
	// configure services (business logic)
	deviceService := domain.NewDeviceService(repo)

	// configure the http server
	server := api.NewServer(ListenAddress)
	// create, configure and assign handlers to routes
	server = server.WithHandler("/api/v0/health/", api.NewHealthHandler())
	server = server.WithHandler("/api/v0/devices/", api.NewDeviceAPIHandler(deviceService))
	// start the server
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
