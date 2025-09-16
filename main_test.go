package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"

	tests := []struct {
		count int
		want  int
	}{
		{0, 0},

		{1, 1},

		{2, 2},

		{100, len(cafeList[city])},

		{len(cafeList[city]) - 1, len(cafeList[city]) - 1},

		{25, len(cafeList[city])},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("count=%d", test.count), func(t *testing.T) {
			requestsUrl := fmt.Sprintf("/cafe?count=%d", test.count)
			response := httptest.NewRecorder()
			req := httptest.NewRequest("GET", requestsUrl, nil)
			handler.ServeHTTP(response, req)
			require.Equal(t, http.StatusOK, response.Code, "request %s must be successful. Status:%d", requestsUrl, response.Code)

			responseBody := strings.TrimSpace(response.Body.String())
			var actualCafes []string
			if responseBody != "" {
				actualCafes = strings.Split(responseBody, ",")
			} else {
				actualCafes = []string{}
			}
			assert.Len(t, actualCafes, test.want, "for request %s: Expected %d cafes, received %d. Response body: '%s'", requestsUrl, test.want, test.count, responseBody)
			if test.count == 0 {
				assert.Equal(t, "", responseBody, "when count=0, the response should be an empty string")
			}

		})
	}
	t.Run("default_count", func(t *testing.T) {
		requestsUrl := fmt.Sprintf("/cafe?city=%s", city)
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", requestsUrl, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code, "request %s must be successful. Status:%d", requestsUrl, response.Code)

		responseBody := strings.TrimSpace(response.Body.String())
		actualCafes := strings.Split(responseBody, ",")

		expectedCount := min(25, len(cafeList[city]))
		assert.Len(t, actualCafes, expectedCount, "for request %s without count: Expected %d cafes, received %d. Response body: '%s'", requestsUrl, expectedCount, len(actualCafes), responseBody)
	})
}
