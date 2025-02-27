// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

type MockServer struct {
	server       *httptest.Server
	applications []types.App
	mu           sync.RWMutex
}

func NewMockServer() *MockServer {
	s := &MockServer{
		applications: make([]types.App, 0),
	}
	s.server = httptest.NewServer(http.HandlerFunc(s.handleRequest))
	return s
}

func (s *MockServer) URL() string {
	return s.server.URL
}

func (s *MockServer) Close() {
	s.server.Close()
}

func (s *MockServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/applications"):
		s.handleCreateApplication(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/applications":
		s.handleListApplications(w)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/applications/"):
		s.handleGetApplication(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/applications/"):
		s.handleDeleteApplication(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

// handleCreateApplication handles application creation requests
func (s *MockServer) handleCreateApplication(w http.ResponseWriter, r *http.Request) {
	var app types.App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	s.applications = append(s.applications, app)
	s.mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(app); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *MockServer) handleListApplications(w http.ResponseWriter) {
	s.mu.RLock()
	if err := json.NewEncoder(w).Encode(s.applications); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()
}

func (s *MockServer) handleGetApplication(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.Path[len("/applications/"):]
	s.mu.RLock()
	for _, app := range s.applications {
		if app.Name == appName {
			if err := json.NewEncoder(w).Encode(app); err != nil {
				http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
				s.mu.RUnlock()
				return
			}
			s.mu.RUnlock()
			return
		}
	}
	s.mu.RUnlock()
	http.Error(w, "Application not found", http.StatusNotFound)
}

func (s *MockServer) handleDeleteApplication(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.Path[len("/applications/"):]
	s.mu.Lock()
	for i, app := range s.applications {
		if app.Name == appName {
			s.applications = append(s.applications[:i], s.applications[i+1:]...)
			s.mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mu.Unlock()
	http.Error(w, "Application not found", http.StatusNotFound)
}
