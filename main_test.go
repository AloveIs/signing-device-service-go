package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/common"
)

func TestErrorMessages(t *testing.T) {

}

// Perform end-to-end testing performing a sequence of requests mocking the
// user behaviour. (similar to acceptance testing)
func TestMain(t *testing.T) {
	// spin-up an instance of the server
	// TODO: here the repository should be the real database (I usually use testcontainers)
	server := configureServer()
	go server.Run()
	N := 1
	wg := &sync.WaitGroup{}
	// emulate N users accessing the API
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			testNonExistingUrl(t)
			testListDevices(t, 0)
			deviceA := testCreateDevice(t)
			deviceB := testCreateDevice(t)
			testListDevices(t, 2)
			testListSignatures(t, 0)
			testSignMessage(t, deviceA)
			testSignMessage(t, deviceB)
			testListSignatures(t, 2)
			// TODO: test retrieve
			// TODO: test error messages
			wg.Done()
		}()
	}
	wg.Wait()
}

func testListDevices(t *testing.T, expected int) {
	resp, err := http.Get("http://localhost:8080/api/v0/devices")
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response body
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	// Decode the response body
	var response struct {
		Data []common.Device `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	devices := response.Data

	if len(devices) != expected {
		t.Errorf("Expected %v devices, got %v", expected, len(devices))
	}
}

func testSignMessage(t *testing.T, device common.Device) common.Signature {
	isb64 := false
	message := "message"
	inputValues := api.SignMessageRequest{
		Message:  &message,
		IsBase64: &isb64,
	}
	jsonValue, err := json.Marshal(inputValues)
	if err != nil {
		t.Errorf("Failed to marshal values: %v", err)
	}
	resp, err := http.Post("http://localhost:8080/api/v0/devices/"+device.ID+"/sign", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Errorf("Create device failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.StatusCode)
	}

	// Decode the response body
	var response struct {
		Data common.Signature `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	signature := response.Data
	return signature
}

// Get the list of signatures and test if the expected length is correct
// TODO: add test to check the IDs match the created one
func testListSignatures(t *testing.T, expected int) {
	resp, err := http.Get("http://localhost:8080/api/v0/signatures")
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response body
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	// Decode the response body
	var response struct {
		Data []common.Signature `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	signatures := response.Data

	if len(signatures) != expected {
		t.Errorf("Expected %v signatures, got %v", expected, len(signatures))
	}
}

func testCreateDevice(t *testing.T) common.Device {
	inputValues := api.CreateDeviceRequest{
		Label:     nil,
		Algorithm: "RSA",
	}
	jsonValue, err := json.Marshal(inputValues)
	if err != nil {
		t.Errorf("Failed to marshal values: %v", err)
	}
	resp, err := http.Post("http://localhost:8080/api/v0/devices/", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Errorf("Create device failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.StatusCode)
	}
	// Decode the response body
	var response struct {
		Data common.Device `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	device := response.Data

	return device
}

func testNonExistingUrl(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/nonexisting")
	if err != nil {
		t.Errorf("Non existing URL check failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func testHealthService(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v0/health")
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}
