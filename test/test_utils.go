package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestRequest represents a request for testing
type TestRequest struct {
	Method  string
	URL     string
	Body    interface{}
	Headers map[string]string
}

// MakeTestRequest performs an HTTP request for testing
func MakeTestRequest(t *testing.T, app *fiber.App, req TestRequest) *http.Response {
	// Create request body if provided
	var reqBody io.Reader
	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		assert.NoError(t, err, "Failed to marshal request body")
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, req.URL, reqBody)
	assert.NoError(t, err, "Failed to create request")

	// Set content type for JSON requests
	httpReq.Header.Set("Content-Type", "application/json")

	// Add additional headers if specified
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Perform request
	resp, err := app.Test(httpReq)
	assert.NoError(t, err, "Request failed")

	return resp
}

// ParseResponse parses the response body into the provided structure
func ParseResponse(t *testing.T, resp *http.Response, target interface{}) {
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	err = json.Unmarshal(body, target)
	assert.NoError(t, err, "Failed to unmarshal response body")
}

// SetupTestApp sets up a Fiber app for testing
func SetupTestApp() *fiber.App {
	app := fiber.New()
	// Add any middleware or routes needed for testing
	return app
}

// WithAuth adds an authentication header to the request
func WithAuth(req TestRequest, token string) TestRequest {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	req.Headers["Authorization"] = "Bearer " + token
	return req
}
