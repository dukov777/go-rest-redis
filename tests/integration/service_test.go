package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceWithTestify(t *testing.T) {
	assert := assert.New(t)

	key := "test-key"
	// Get the current time
	valueNow := time.Now().Format("2006-01-02 15:04:05")

	// Create JSON payload
	data := map[string]string{
		"key":   key,
		"value": valueNow,
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

	redisValue, err := getRedisValue(key)
	assert.NoError(err, "Getting Redis value should not error")
	assert.Equal(valueNow, redisValue)

}

// getRedisValue uses redis-cli to get a value for a given key.
func getRedisValue(key string) (string, error) {
	cmd := exec.Command("redis-cli", "GET", key)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	// Trim space to clean up the output
	return strings.TrimSpace(string(output)), nil
}
