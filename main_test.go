package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/AloveIs/signing-device-service-go/api"
	"github.com/AloveIs/signing-device-service-go/common"
)

// Perform end-to-end testing performing a sequence of requests mocking the
// user behaviour. (similar to acceptance testing)
func TestMain(t *testing.T) {
	// spin-up an instance of the server
	// TODO: here the repository should be the real database (I usually use testcontainers)
	server := configureServer()
	go func() {
		if err := server.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	N := 1
	wg := &sync.WaitGroup{}
	// emulate N users accessing the API
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			// TODO: this test needs should add more test cases
			testNonExistingUrl(t)
			testListDevices(t, []common.Device{})
			deviceA := testCreateDevice(t)
			deviceB := testCreateDevice(t)
			testListDevices(t, []common.Device{deviceA, deviceB})
			testListSignatures(t, []common.Signature{})
			sigA := testSignMessage(t, deviceA)
			sigB := testSignMessage(t, deviceB)
			testListSignatures(t, []common.Signature{sigA, sigB})
			// test retrieve
			testRetrieveSignature(t, sigA)
			testRetrieveDevice(t, deviceA)
			// test error messages
			testRetrieveDeviceFailure(t, "IMPOSSIBLE_DEVICE_ID")
			testRetrieveSignatureFailure(t, "IMPOSSIBLE_DEVICE_ID")
			// test invalid payloads
			testSignatureFailure(t, deviceA)
			testDeviceCreationFailure(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func testListDevices(t *testing.T, expected []common.Device) {
	resp, err := http.Get("http://localhost:8080/api/v0/devices/")
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response statusCode
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

	if len(devices) != len(expected) {
		t.Errorf("Expected %v devices, got %v", len(expected), len(devices))
	}

	idSet := make(map[string]struct{})
	for _, device := range expected {
		idSet[device.ID] = struct{}{}
	}

	for _, device := range devices {
		if _, ok := idSet[device.ID]; !ok {
			t.Errorf("Unexpected device ID: %v", device.ID)
		}
	}
}

func testRetrieveDevice(t *testing.T, expected common.Device) {
	resp, err := http.Get("http://localhost:8080/api/v0/devices/" + expected.ID)
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response statusCode
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Decode the response body
	var response struct {
		Data common.Device `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	device := response.Data

	if device != expected {
		t.Errorf("Expected %v ==  %v", device, expected)
	}
}

func testRetrieveSignature(t *testing.T, expected common.Signature) {
	resp, err := http.Get("http://localhost:8080/api/v0/signatures/" + expected.ID)
	if err != nil {
		t.Errorf("Retrieve signature failed: %v", err)
	}
	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
	// Decode the response body
	var response struct {
		Data common.Signature `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}
	signature := response.Data

	if signature != expected {
		t.Errorf("Expected %v ==  %v", signature, expected)
	}

}

func testRetrieveDeviceFailure(t *testing.T, deviceId string) {
	resp, err := http.Get("http://localhost:8080/api/v0/devices/" + deviceId)
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response statusCode
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.StatusCode)
	}
}

func testRetrieveSignatureFailure(t *testing.T, signatureId string) {
	resp, err := http.Get("http://localhost:8080/api/v0/signatures/" + signatureId)
	if err != nil {
		t.Errorf("List devices failed: %v", err)
	}
	// Check the response statusCode
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.StatusCode)
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
func testListSignatures(t *testing.T, expected []common.Signature) {
	resp, err := http.Get("http://localhost:8080/api/v0/signatures/")
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

	if len(signatures) != len(expected) {
		t.Errorf("Expected %v signatures, got %v", len(expected), len(signatures))
	}

	idSet := make(map[string]struct{})
	for _, signature := range expected {
		idSet[signature.ID] = struct{}{}
	}

	for _, signature := range signatures {
		if _, ok := idSet[signature.ID]; !ok {
			t.Errorf("Unexpected signature ID: %v", signature.ID)
		}
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

func testSignatureFailure(t *testing.T, device common.Device) {
	IMPOSSIBLE_DEVICE_ID := "IMPOSSIBLE_DEVICE_ID"
	message := "message"
	isb64 := false

	testCases := []struct {
		jsonMessage   string
		expctedStatus int
	}{
		{"-", http.StatusBadRequest},

		{"{}", http.StatusUnprocessableEntity},
		{`{"message": null}`, http.StatusUnprocessableEntity},
		{`{"message": null, "isBase64": true}`, http.StatusUnprocessableEntity},
		{`{"message": "abc"}`, http.StatusUnprocessableEntity},
	}

	// run test cases
	for _, tc := range testCases {

		// test impossible device id, expected not found
		resp, err := http.Post("http://localhost:8080/api/v0/devices/"+device.ID+"/sign", "application/json", bytes.NewBuffer([]byte(tc.jsonMessage)))
		if err != nil {
			t.Errorf("Create device failed: %v", err)
		}
		if resp.StatusCode != tc.expctedStatus {
			t.Errorf("Expected status %d, got %v", tc.expctedStatus, resp.StatusCode)
		}

	}
	inputValues := api.SignMessageRequest{
		Message:  &message,
		IsBase64: &isb64,
	}
	jsonValue, err := json.Marshal(inputValues)
	if err != nil {
		t.Errorf("Failed to marshal values: %v", err)
	}

	// test impossible device id, expected not found
	resp, err := http.Post("http://localhost:8080/api/v0/devices/"+IMPOSSIBLE_DEVICE_ID+"/sign", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Errorf("Create device failed: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Created, got %v", resp.StatusCode)
	}
}

func testDeviceCreationFailure(t *testing.T) {
	testCases := []struct {
		jsonMessage   string
		expctedStatus int
	}{
		{"-", http.StatusBadRequest},
		{"{}", http.StatusUnprocessableEntity},
		{`{"algorithm": null}`, http.StatusUnprocessableEntity},
		{`{"algorithm": "AES", "label": "some-label"}`, http.StatusUnprocessableEntity},
		{`{"algorithm": "AES", "label": null}`, http.StatusUnprocessableEntity},
	}

	// run test cases
	for _, tc := range testCases {

		// test impossible device id, expected not found
		resp, err := http.Post("http://localhost:8080/api/v0/devices/", "application/json", bytes.NewBuffer([]byte(tc.jsonMessage)))
		if err != nil {
			t.Errorf("Create device failed: %v", err)
		}
		if resp.StatusCode != tc.expctedStatus {
			t.Errorf("Expected status %d, got %v", tc.expctedStatus, resp.StatusCode)
		}

	}
}
