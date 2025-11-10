// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/armis"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestAppendAPIError_NilError tests appendAPIError with nil error input.
func TestAppendAPIError_NilError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		title             string
		err               error
		expectedSummary   string
		expectedDetail    string
		expectedDiagCount int
	}{
		{
			name:              "nil error should add default message",
			title:             "Test Error",
			err:               nil,
			expectedSummary:   "Test Error",
			expectedDetail:    "(no additional error details returned)",
			expectedDiagCount: 1,
		},
		{
			name:              "nil error with different title",
			title:             "Failed to create resource",
			err:               nil,
			expectedSummary:   "Failed to create resource",
			expectedDetail:    "(no additional error details returned)",
			expectedDiagCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var diags diag.Diagnostics
			appendAPIError(&diags, tt.title, tt.err)

			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiagCount, len(diags))
			}

			if len(diags) > 0 {
				if diags[0].Summary() != tt.expectedSummary {
					t.Errorf("Expected summary '%s', got '%s'", tt.expectedSummary, diags[0].Summary())
				}
				if diags[0].Detail() != tt.expectedDetail {
					t.Errorf("Expected detail '%s', got '%s'", tt.expectedDetail, diags[0].Detail())
				}
			}
		})
	}
}

// TestAppendAPIError_APIError tests appendAPIError with armis.APIError type.
func TestAppendAPIError_APIError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		title             string
		statusCode        int
		body              []byte
		validateDetail    func(t *testing.T, detail string)
		expectedDiagCount int
	}{
		{
			name:       "400 Bad Request with JSON body",
			title:      "Request Failed",
			statusCode: 400,
			body:       []byte(`{"error":"Invalid request","code":"BAD_REQUEST"}`),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 400 Bad Request") {
					t.Errorf("Expected detail to contain 'API error 400 Bad Request', got: %s", detail)
				}
				if !strings.Contains(detail, `"error": "Invalid request"`) {
					t.Errorf("Expected pretty-printed JSON in detail, got: %s", detail)
				}
				if !strings.Contains(detail, "Response body:") {
					t.Errorf("Expected 'Response body:' in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "401 Unauthorized with JSON body",
			title:      "Authentication Failed",
			statusCode: 401,
			body:       []byte(`{"message":"Unauthorized"}`),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 401 Unauthorized") {
					t.Errorf("Expected detail to contain 'API error 401 Unauthorized', got: %s", detail)
				}
				if !strings.Contains(detail, `"message": "Unauthorized"`) {
					t.Errorf("Expected pretty-printed JSON in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "403 Forbidden with plain text body",
			title:      "Access Denied",
			statusCode: 403,
			body:       []byte("Access forbidden"),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 403 Forbidden") {
					t.Errorf("Expected detail to contain 'API error 403 Forbidden', got: %s", detail)
				}
				if !strings.Contains(detail, "Access forbidden") {
					t.Errorf("Expected plain text body in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "404 Not Found with JSON body",
			title:      "Resource Not Found",
			statusCode: 404,
			body:       []byte(`{"error":"Resource not found","id":"12345"}`),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 404 Not Found") {
					t.Errorf("Expected detail to contain 'API error 404 Not Found', got: %s", detail)
				}
				if !strings.Contains(detail, `"error": "Resource not found"`) {
					t.Errorf("Expected pretty-printed JSON in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "500 Internal Server Error with JSON body",
			title:      "Server Error",
			statusCode: 500,
			body:       []byte(`{"error":"Internal server error","trace":"stack trace here"}`),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 500 Internal Server Error") {
					t.Errorf("Expected detail to contain 'API error 500 Internal Server Error', got: %s", detail)
				}
				if !strings.Contains(detail, `"error": "Internal server error"`) {
					t.Errorf("Expected pretty-printed JSON in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "503 Service Unavailable with HTML body",
			title:      "Service Unavailable",
			statusCode: 503,
			body:       []byte("<html><body>Service Unavailable</body></html>"),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 503 Service Unavailable") {
					t.Errorf("Expected detail to contain 'API error 503 Service Unavailable', got: %s", detail)
				}
				if !strings.Contains(detail, "<html><body>Service Unavailable</body></html>") {
					t.Errorf("Expected HTML body in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "empty response body",
			title:      "Empty Response",
			statusCode: 500,
			body:       []byte{},
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 500 Internal Server Error") {
					t.Errorf("Expected detail to contain 'API error 500 Internal Server Error', got: %s", detail)
				}
				if !strings.Contains(detail, "(empty response body)") {
					t.Errorf("Expected '(empty response body)' in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "whitespace-only response body",
			title:      "Whitespace Response",
			statusCode: 400,
			body:       []byte("   \n\t   "),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 400 Bad Request") {
					t.Errorf("Expected detail to contain 'API error 400 Bad Request', got: %s", detail)
				}
				if !strings.Contains(detail, "(empty response body)") {
					t.Errorf("Expected '(empty response body)' for whitespace, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:       "complex nested JSON body",
			title:      "Complex Error",
			statusCode: 422,
			body:       []byte(`{"errors":[{"field":"name","message":"required"},{"field":"email","message":"invalid format"}],"code":"VALIDATION_ERROR"}`),
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error 422") {
					t.Errorf("Expected detail to contain 'API error 422', got: %s", detail)
				}
				// Verify JSON is pretty-printed (has indentation)
				if !strings.Contains(detail, `  "errors": [`) {
					t.Errorf("Expected indented JSON in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var diags diag.Diagnostics
			apiErr := &armis.APIError{
				StatusCode: tt.statusCode,
				Body:       tt.body,
			}

			appendAPIError(&diags, tt.title, apiErr)

			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiagCount, len(diags))
			}

			if len(diags) > 0 {
				if diags[0].Summary() != tt.title {
					t.Errorf("Expected summary '%s', got '%s'", tt.title, diags[0].Summary())
				}
				tt.validateDetail(t, diags[0].Detail())
			}
		})
	}
}

// TestAppendAPIError_GenericError tests appendAPIError with generic Go error types.
func TestAppendAPIError_GenericError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		title             string
		err               error
		expectedSummary   string
		validateDetail    func(t *testing.T, detail string)
		expectedDiagCount int
	}{
		{
			name:            "simple error from errors.New",
			title:           "Operation Failed",
			err:             errors.New("connection timeout"), //nolint:err113 // test fixture
			expectedSummary: "Operation Failed",
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error: connection timeout") {
					t.Errorf("Expected detail to contain 'API error: connection timeout', got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:            "formatted error from fmt.Errorf",
			title:           "Configuration Error",
			err:             fmt.Errorf("invalid configuration: %s", "missing required field"), //nolint:err113 // test fixture
			expectedSummary: "Configuration Error",
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "API error: invalid configuration: missing required field") {
					t.Errorf("Expected formatted error message in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:            "wrapped error",
			title:           "Database Error",
			err:             fmt.Errorf("failed to connect: %w", errors.New("host unreachable")), //nolint:err113 // test fixture
			expectedSummary: "Database Error",
			validateDetail: func(t *testing.T, detail string) {
				expected := "API error: failed to connect: host unreachable"
				if !strings.Contains(detail, expected) {
					t.Errorf("Expected detail to contain '%s', got: %s", expected, detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:            "error with special characters",
			title:           "Parse Error",
			err:             errors.New("invalid JSON: unexpected character '}'"), //nolint:err113 // test fixture
			expectedSummary: "Parse Error",
			validateDetail: func(t *testing.T, detail string) {
				if !strings.Contains(detail, "invalid JSON: unexpected character '}'") {
					t.Errorf("Expected special characters preserved in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var diags diag.Diagnostics
			appendAPIError(&diags, tt.title, tt.err)

			if len(diags) != tt.expectedDiagCount {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiagCount, len(diags))
			}

			if len(diags) > 0 {
				if diags[0].Summary() != tt.expectedSummary {
					t.Errorf("Expected summary '%s', got '%s'", tt.expectedSummary, diags[0].Summary())
				}
				tt.validateDetail(t, diags[0].Detail())
			}
		})
	}
}

// TestAppendAPIError_EdgeCases tests edge cases for appendAPIError.
func TestAppendAPIError_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		title             string
		err               error
		validateDiag      func(t *testing.T, diags diag.Diagnostics)
		expectedDiagCount int
	}{
		{
			name:  "empty title string",
			title: "",
			err:   errors.New("some error"), //nolint:err113 // test fixture
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				if diags[0].Summary() != "" {
					t.Errorf("Expected empty summary, got '%s'", diags[0].Summary())
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "very long error message",
			title: "Long Error",
			err:   errors.New(strings.Repeat("x", 1000)), //nolint:err113 // test fixture
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				if !strings.Contains(diags[0].Detail(), strings.Repeat("x", 1000)) {
					t.Error("Expected very long error message to be preserved")
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "malformed JSON in APIError body",
			title: "Malformed JSON",
			err: &armis.APIError{
				StatusCode: 500,
				Body:       []byte(`{"incomplete": "json"`),
			},
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				// Should still include the body even if JSON is invalid (not pretty-printed)
				if !strings.Contains(diags[0].Detail(), `{"incomplete": "json"`) {
					t.Errorf("Expected malformed JSON to be included as-is, got: %s", diags[0].Detail())
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "APIError with newlines and tabs in body",
			title: "Special Characters",
			err: &armis.APIError{
				StatusCode: 400,
				Body:       []byte("Error message\nwith\nnewlines\tand\ttabs"),
			},
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				detail := diags[0].Detail()
				if !strings.Contains(detail, "Error message") {
					t.Errorf("Expected error message in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "APIError with Unicode characters",
			title: "Unicode Test",
			err: &armis.APIError{
				StatusCode: 400,
				Body:       []byte(`{"message":"Error: 你好 🚀 Ñoño"}`),
			},
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				detail := diags[0].Detail()
				if !strings.Contains(detail, "你好") || !strings.Contains(detail, "🚀") || !strings.Contains(detail, "Ñoño") {
					t.Errorf("Expected Unicode characters preserved in detail, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "APIError with only spaces in body",
			title: "Spaces Only",
			err: &armis.APIError{
				StatusCode: 500,
				Body:       []byte("     "),
			},
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				if !strings.Contains(diags[0].Detail(), "(empty response body)") {
					t.Errorf("Expected '(empty response body)' for spaces-only body, got: %s", diags[0].Detail())
				}
			},
			expectedDiagCount: 1,
		},
		{
			name:  "error with newlines",
			title: "Multi-line Error",
			err:   errors.New("line 1\nline 2\nline 3"), //nolint:err113 // test fixture
			validateDiag: func(t *testing.T, diags diag.Diagnostics) {
				if len(diags) != 1 {
					t.Errorf("Expected 1 diagnostic, got %d", len(diags))
					return
				}
				detail := diags[0].Detail()
				if !strings.Contains(detail, "line 1\nline 2\nline 3") {
					t.Errorf("Expected multi-line error preserved, got: %s", detail)
				}
			},
			expectedDiagCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var diags diag.Diagnostics
			appendAPIError(&diags, tt.title, tt.err)
			tt.validateDiag(t, diags)
		})
	}
}
