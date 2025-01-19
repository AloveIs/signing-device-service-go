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
	server := configureServer()
	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}

func configureServer() *api.Server {
	// create the repositories (database)
	deviceRepo := persistence.NewInMemoryDeviceDb()
	signatureRepo := persistence.NewInMemorySignatureDb()

	// configure services (business logic)
	deviceService := domain.NewDeviceService(deviceRepo, signatureRepo)
	signatureService := domain.NewSignatureService(signatureRepo)

	// configure the http server
	server := api.NewServer(ListenAddress)

	// create, configure and assign handlers to routes
	server = server.WithHandler("/api/v0/health/", api.NewHealthHandler())
	server = server.WithHandler("/api/v0/devices/", api.NewDeviceAPIHandler(deviceService))
	server = server.WithHandler("/api/v0/signatures/", api.NewSignatureAPIHandler(signatureService))

	// start the server
	return server
}
