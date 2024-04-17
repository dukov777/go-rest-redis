package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceWithTestify(t *testing.T) {
	assert := assert.New(t)

	// Create JSON payload
	data := map[string]string{
		"key":   "12356",
		"value": time.Now().Format("2006-01-02 15:04:05"),
	}
	jsonData, err := json.Marshal(data)
	assert.NoError(err, "Encoding JSON should not error")

	// Make the POST request
	resp, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(err, "POST request should not error")
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(http.StatusOK, resp.StatusCode, "Expected status code 200")

	// Read and check the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err, "Reading response body should not error")
	expectedBody := "All key-value pairs received and stored in Redis.\n"
	assert.Equal(expectedBody, string(responseBody), "Response body should match expected content")
}
