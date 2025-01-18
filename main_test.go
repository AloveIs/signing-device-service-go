package main

import (
	"net/http"
	"testing"
)

func testDeviceEndpoints(t *testing.T) {
	// Test GET devices endpoint
	resp, err := http.Get("http://localhost:8080/devices")
	if err != nil {
		t.Errorf("GET devices failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for GET devices, got %v", resp.StatusCode)
	}

	// Test GET single device endpoint
	resp, err = http.Get("http://localhost:8080/devices/1")
	if err != nil {
		t.Errorf("GET single device failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for GET single device, got %v", resp.StatusCode)
	}

	// Test POST new device
	req, _ := http.NewRequest("POST", "http://localhost:8080/devices", nil)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("POST device failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created for POST device, got %v", resp.StatusCode)
	}

	// Test PUT update device
	req, _ = http.NewRequest("PUT", "http://localhost:8080/devices/1", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("PUT device failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for PUT device, got %v", resp.StatusCode)
	}

	// Test DELETE device
	req, _ = http.NewRequest("DELETE", "http://localhost:8080/devices/1", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("DELETE device failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK for DELETE device, got %v", resp.StatusCode)
	}
}

func testNonExistingUrl(t *testing.T) {
	// Test server is running
	resp, err := http.Get("http://localhost:8080/nonexisting")
	if err != nil {
		t.Errorf("Non existing URL check failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func testHealthService(t *testing.T) {
	// Test server is running
	resp, err := http.Get("http://localhost:8080/api/v0/health")
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func TestHumanReadableErrorMessages(t *testing.T) {

}

func testGenerateKey() {}
func testGetKey()      {}
func testSign()        {}

// Perform end-to-end testing
func TestMain(t *testing.T) {
	// spin-up an instance of the server
	// TODO: here the repository should be the real database (I usually use testcontainers)
	server := configureServer()
	go server.Run()
	N := 100
	// emulate N users accessing the API
	for i := 0; i < N; i++ {
		go func() {
			testHealthService(t)
			testNonExistingUrl(t)
			testGenerateKey()
			testGetKey()
			testSign()
		}()
	}
}
