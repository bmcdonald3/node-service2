package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/node-service/internal/api"
)

func TestGetComposedNodeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/composed/nodes/x1000c0s1b0n0", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.GetComposedNodeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["xname"] != "x1000c0s1b0n0" {
		t.Errorf("Expected xname x1000c0s1b0n0, got %v", response["xname"])
	}
	if response["effectiveProfile"] != "compute-new" {
		t.Errorf("Expected effectiveProfile compute-new, got %v", response["effectiveProfile"])
	}
}