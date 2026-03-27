package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func TestProfileBindingAPI(t *testing.T) {
	t.Run("CreateInvalidBinding", func(t *testing.T) {
		payload := []byte(`{"apiVersion":"v1","kind":"ProfileBinding","metadata":{"name":"test-binding-invalid"},"spec":{"targetKind":"Node","targetName":"x1000c0s1b0n0"}}`)
		resp, err := http.Post(baseURL+"/profilebindings", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("CreateValidBinding", func(t *testing.T) {
		payload := []byte(`{"apiVersion":"v1","kind":"ProfileBinding","metadata":{"name":"test-binding-valid"},"spec":{"targetKind":"Node","targetName":"x1000c0s1b0n0","profile":"compute-new"}}`)
		resp, err := http.Post(baseURL+"/profilebindings", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 201 or 200, got %d", resp.StatusCode)
		}
	})

	t.Run("GetBindings", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/profilebindings")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var bindings []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&bindings); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(bindings) == 0 {
			t.Errorf("Expected at least one binding, got 0")
		}
	})
}
