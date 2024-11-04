package handlers

import (
	"bytes"
	"encoding/json"
	"meter_flow/server"
	"meter_flow/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScheduleCalls(t *testing.T) {
	storage := storage.NewDummyStorage()
	server := server.NewServer(storage)

	// Register a test resource
	registerTestResource(t, server)

	// Test cases
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedOutput string
	}{
		{
			name:           "Valid schedule",
			requestBody:    `{"resource_name":"test_resource", "num_calls":5}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Resource not found",
			requestBody:    `{"resource_name":"non_existent_resource", "num_calls":5}`,
			expectedStatus: http.StatusNotFound,
			expectedOutput: "Resource not found\n",
		},
		{
			name:           "Invalid request",
			requestBody:    `{"resource_name":"test_resource", "num_calls":-1}`,
			expectedStatus: http.StatusBadRequest,
			expectedOutput: "Invalid request\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("POST", "/schedule", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Errorf("failed to create request: %v", err)
			}

			// Create a new HTTP recorder
			rr := httptest.NewRecorder()

			// Call the scheduleCalls handler
			handler := ScheduleCalls(server)
			handler(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check the response body if expected
			if tc.expectedOutput != "" && rr.Body.String() != tc.expectedOutput {
				t.Errorf("expected response body %q, got %q", tc.expectedOutput, rr.Body.String())
			}
		})
	}
}

func registerTestResource(t *testing.T, server *server.Server) {
	// Register the "test_resource"
	resourceData := struct {
		Name         string `json:"name"`
		RequestCount int    `json:"request_count"`
		TimeFrame    int    `json:"time_frame"`
	}{
		Name:         "test_resource",
		RequestCount: 10,
		TimeFrame:    60,
	}

	requestBody, err := json.Marshal(resourceData)
	if err != nil {
		t.Errorf("failed to marshal resource data: %v", err)
	}

	req, err := http.NewRequest("POST", "/resources", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := RegisterResource(server)
	handler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
}
