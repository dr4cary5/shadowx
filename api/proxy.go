package handler

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

type Response struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	GoVersion string `json:"go_version"`
	Region    string `json:"region"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status:    "ok",
		Message:   "ShadowX serverless function is running",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		GoVersion: runtime.Version(),
		Region:    r.Header.Get("X-Vercel-Edge-Region"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
