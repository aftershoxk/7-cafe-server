package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(100, len(cafeList["tula"]))},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=tula&count=%d", v.count), nil)

		handler.ServeHTTP(response, req)

		body := strings.TrimSpace(response.Body.String())

		if v.want == 0 {
			assert.Equal(t, "", strings.TrimSpace(body))
		} else {
			cafes := strings.Split(body, ",")
			assert.Equal(t, v.want, len(cafes))
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)

		handler.ServeHTTP(response, req)

		body := strings.TrimSpace(response.Body.String())

		if v.wantCount == 0 {
			assert.Equal(t, "", body)
		} else {
			cafes := strings.Split(body, ",")
			assert.Equal(t, v.wantCount, len(cafes))
			for _, cafe := range cafes {
				assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)))
			}
		}
	}
}
