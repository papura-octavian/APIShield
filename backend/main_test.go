package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkParameter(t *testing.T, got, expected, context string) {
	t.Helper()

	if got != expected {
		t.Errorf("%q: expected %q, got %q", context, expected, got)
	}
}

func TestHelloHandlerPost(t *testing.T) {
	req := httptest.NewRequest("POST", "/users/42", nil)
	rec := httptest.NewRecorder()

	helloHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got Response
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON:%v", err)
	}

	checkParameter(t, got.Method, "POST", "Method")
	checkParameter(t, got.Message, "hello from backend!", "Message")
	checkParameter(t, got.Path, "/users/42", "Path")
}

func TestHelloHandlerGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	helloHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got Response
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("invalid JSON:%v", err)
	}

	checkParameter(t, got.Method, "GET", "Method")
	checkParameter(t, got.Message, "hello from backend!", "Message")
	checkParameter(t, got.Path, "/", "Path")
}
