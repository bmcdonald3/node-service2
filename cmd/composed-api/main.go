package main

import (
	"log"
	"net/http"

	"github.com/user/node-service/internal/api"
)

func main() {
	// Mount the custom Composed API handler
	http.HandleFunc("/composed/nodes/", api.GetComposedNodeHandler)

	// Mock the metadata and boot service binding endpoints
	http.HandleFunc("/bindings/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Sidecar: Composed API and Downstream Mocks listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}