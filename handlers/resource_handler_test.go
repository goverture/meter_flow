package handlers

import (
	"bytes"
	"encoding/json"
	"meter_flow/model"
	"meter_flow/server"
	"meter_flow/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterResource(t *testing.T) {
	// Create a new server
	storage := storage.NewDummyStorage()
	server := server.NewServer(storage)

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
			handler := RegisterResource(server)
			handler(rr, req)

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

func TestListResources(t *testing.T) {
	storage := storage.NewDummyStorage()
	server := server.NewServer(storage)

	// Register some test resources
	server.Resources = map[string]model.Resource{
		"test_resource_1": {
			Name:         "test_resource_1",
			RequestCount: 10,
			TimeFrame:    60,
		},
		"test_resource_2": {
			Name:         "test_resource_2",
			RequestCount: 20,
			TimeFrame:    120,
		},
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/resources", nil)
	if err != nil {
		t.Errorf("failed to create request: %v", err)
	}

	// Create a new HTTP recorder
	rr := httptest.NewRecorder()

	// Call the listResources handler
	handler := ListResources(server)
	handler(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var resources []model.Resource
	if err := json.NewDecoder(rr.Body).Decode(&resources); err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(resources))
	}
}

func TestUpdateResource(t *testing.T) {
	storage := storage.NewDummyStorage()
	server := server.NewServer(storage)

	// Register a test resource
	server.Resources = map[string]model.Resource{
		"test_resource": {
			Name:         "test_resource",
			RequestCount: 10,
			TimeFrame:    60,
		},
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedOutput string
	}{
		{
			name:           "Valid update",
			requestBody:    `{"name":"test_resource", "request_count":20, "time_frame":120}`,
			expectedStatus: http.StatusOK,
			expectedOutput: "Resource test_resource updated with limit of 20 requests per 120 seconds\n",
		},
		{
			name:           "Resource not found",
			requestBody:    `{"name":"non_existent_resource", "request_count":20, "time_frame":120}`,
			expectedStatus: http.StatusNotFound,
			expectedOutput: "Resource not found\n",
		},
		{
			name:           "Invalid request",
			requestBody:    `{"name":"test_resource", "request_count":-1, "time_frame":120}`,
			expectedStatus: http.StatusBadRequest,
			expectedOutput: "Invalid request\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("PUT", "/resource", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Errorf("failed to create request: %v", err)
			}

			// Create a new HTTP recorder
			rr := httptest.NewRecorder()

			// Call the updateResource handler
			handler := UpdateResource(server)
			handler(rr, req)

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

func TestDeleteResource(t *testing.T) {
	storage := storage.NewDummyStorage()
	server := server.NewServer(storage)

	// Register a test resource
	server.Resources = map[string]model.Resource{
		"test_resource": {
			Name:         "test_resource",
			RequestCount: 10,
			TimeFrame:    60,
		},
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedOutput string
	}{
		{
			name:           "Valid delete",
			requestBody:    `{"name":"test_resource"}`,
			expectedStatus: http.StatusOK,
			expectedOutput: "Resource test_resource deleted\n",
		},
		{
			name:           "Resource not found",
			requestBody:    `{"name":"non_existent_resource"}`,
			expectedStatus: http.StatusNotFound,
			expectedOutput: "Resource not found\n",
		},
		{
			name:           "Invalid request",
			requestBody:    `{"name":""}`,
			expectedStatus: http.StatusBadRequest,
			expectedOutput: "Invalid request\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new HTTP request
			req, err := http.NewRequest("DELETE", "/resource", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Errorf("failed to create request: %v", err)
			}

			// Create a new HTTP recorder
			rr := httptest.NewRecorder()

			// Call the deleteResource handler
			handler := DeleteResource(server)
			handler(rr, req)

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
