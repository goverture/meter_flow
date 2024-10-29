package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterResource(t *testing.T) {
	// Create a new server
	server := NewServer()

	// Test cases
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedOutput string
	}{
		{
			name:           "Valid registration",
			requestBody:    `{"name":"test_resource", "request_count":10, "time_frame":60}`,
			expectedStatus: http.StatusCreated,
			expectedOutput: "Resource test_resource with limit of 10 requests per 60 seconds registered\n",
		},
		{
			name:           "Duplicate resource",
			requestBody:    `{"name":"test_resource", "request_count":10, "time_frame":60}`,
			expectedStatus: http.StatusConflict,
			expectedOutput: "Resource already exists\n",
		},
		{
			name:           "Invalid request",
			requestBody:    `{"name":"test_resource", "request_count":-1, "time_frame":60}`,
			expectedStatus: http.StatusBadRequest,
			expectedOutput: "Invalid request\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("POST", "/resources", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Errorf("failed to create request: %v", err)
			}

			// Create a new HTTP recorder
			rr := httptest.NewRecorder()

			// Call the registerResource handler
			server.registerResource(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tc.expectedOutput {
				t.Errorf("expected response body %q, got %q", tc.expectedOutput, rr.Body.String())
			}
		})
	}
}
