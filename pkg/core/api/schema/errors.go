// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}
