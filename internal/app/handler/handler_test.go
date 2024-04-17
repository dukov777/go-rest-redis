package handler // replace with your actual package name

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"my-go-project/internal/app/mocks"

	"github.com/golang/mock/gomock"
)

func TestPayloadHandler(t *testing.T) {
	// Mock the Redis client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRedisClient := mocks.NewMockRedisClient(ctrl) // You will need to generate this mock with a tool like mockgen

	// Example of setting up a test for a successful request
	t.Run("success", func(t *testing.T) {

		handler := PayloadHandler(mockRedisClient)
		mockRedisClient.
			EXPECT().
			SetKey(gomock.Any(), "testKey", "testValue", time.Duration(0)).
			Return(nil).
			Times(1)

		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"key":"testKey", "value":"testValue"}`))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})
	t.Run("fail", func(t *testing.T) {

		handler := PayloadHandler(mockRedisClient)
		mockRedisClient.
			EXPECT().
			SetKey(gomock.Any(), "testKey", "testValue", time.Duration(0)).
			Return(errors.New("error")).
			Times(1)

		r := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"key":"testKey", "value":"testValue"}`))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %v", w.Code)
		}
	})
}
