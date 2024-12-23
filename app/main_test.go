package main

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name        string
		rawRequest  string
		expected    *Request
		expectError bool
	}{
		{
			name:       "Valid GET request with User-Agent",
			rawRequest: "GET /user-agent HTTP/1.1\r\nUser-Agent: TestAgent\r\n\r\n",
			expected: &Request{
				Method:      "GET",
				URL:         "/user-agent",
				HTTPVersion: "HTTP/1.1",
				UserAgent:   "TestAgent",
			},
			expectError: false,
		},
		{
			name:        "Request with no User-Agent header",
			rawRequest:  "GET /user-agent HTTP/1.1\r\n\r\n",
			expected:    &Request{Method: "GET", URL: "/user-agent", HTTPVersion: "HTTP/1.1", UserAgent: ""},
			expectError: false,
		},
		{
			name:        "Malformed request line",
			rawRequest:  "INVALID_REQUEST_LINE\r\n\r\n",
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := parseRequest(tt.rawRequest)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if req.Method != tt.expected.Method || req.URL != tt.expected.URL || req.HTTPVersion != tt.expected.HTTPVersion || req.UserAgent != tt.expected.UserAgent {
					t.Errorf("expected %+v, got %+v", tt.expected, req)
				}
			}
		})
	}
}
