package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Mock responses
var mockDeployment = map[string]interface{}{
	"namespace":        "test-namespace",
	"templateID":       "0001",
	"templateName":     "K-d chat",
	"deploymentStatus": "running",
}

func main() {
	// API endpoints
	http.HandleFunc("/startUserDeployment/", handleStartUserDeployment)
	http.HandleFunc("/saveCustomDomain/", handleSaveCustomDomain)
	http.HandleFunc("/getDeployments/", handleGetDeployments)
	http.HandleFunc("/getDeploymentInfo/", handleGetDeploymentInfo)
	http.HandleFunc("/api/v1/deploy", handleDeploy)
	http.HandleFunc("/api/v1/services/configure", handleConfigure)
	http.HandleFunc("/api/v1/services/", handleServiceConnections)
	http.HandleFunc("/api/v1/deployments/", handleScaleDeployment)
	http.HandleFunc("/api/v1/ai/suggest", handleAISuggestions)
	http.HandleFunc("/", handleRoot)

	http.HandleFunc("/api/v1/deployments", handleGetDeployments)
	http.HandleFunc("/api/v1/deployment", handleGetDeploymentInfo)
	http.HandleFunc("/api/v1/services", handleServiceConnections)
	http.HandleFunc("/api/v1/scale", handleScaleDeployment)

	fmt.Println("Starting mock API server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStartUserDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"message":   "Deployment started",
		"namespace": "fantastic-fox",
		"url":       "https://fantastic-fox-my-mern-app.alpha.nexlayer.io",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSaveCustomDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"message": "Domain saved",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetDeployments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"deployments": []interface{}{mockDeployment},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetDeploymentInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"deployment": mockDeployment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleConfigure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleServiceConnections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connections := []map[string]string{
		{
			"from":        "frontend",
			"to":          "backend",
			"description": "HTTP/REST",
		},
		{
			"from":        "backend",
			"to":          "database",
			"description": "PostgreSQL",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connections)
}

func handleScaleDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleAISuggestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	suggestions := map[string][]string{
		"suggestions": {
			"Use a load balancer",
			"Add monitoring",
			"Set up auto-scaling",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Welcome to the mock API server!")
}
