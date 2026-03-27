package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ComposedNodeResponse struct {
	Xname            string `json:"xname"`
	Role             string `json:"role"`
	EffectiveProfile string `json:"effectiveProfile"`
	BootParams       string `json:"bootParams"`
	CloudInit        string `json:"cloudInit"`
}

func GetComposedNodeHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	xname := parts[len(parts)-1]

	response := ComposedNodeResponse{
		Xname:            xname,
		Role:             "compute",
		EffectiveProfile: "compute-new",
		BootParams:       "console=ttyS0 root=/dev/sda1 profile=compute-new",
		CloudInit:        "{ \"user-data\": \"...\" }",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}