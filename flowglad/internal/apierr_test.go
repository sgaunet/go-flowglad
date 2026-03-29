package internal

import (
	"strings"
	"testing"
)

func TestAPIError_Error_WithIssues(t *testing.T) {
	e := &APIError{
		StatusCode: 400,
		Code:       "validation_error",
		Message:    "invalid input",
		Issues:     []string{"email is required", "name is required"},
	}
	got := e.Error()
	if !strings.Contains(got, "400") {
		t.Errorf("expected status 400 in error string, got: %s", got)
	}
	if !strings.Contains(got, "issues:") {
		t.Errorf("expected 'issues:' in error string, got: %s", got)
	}
}

func TestAPIError_Error_WithoutIssues(t *testing.T) {
	e := &APIError{
		StatusCode: 401,
		Code:       "unauthorized",
		Message:    "invalid API key",
	}
	got := e.Error()
	if !strings.Contains(got, "401") {
		t.Errorf("expected status 401 in error string, got: %s", got)
	}
	if strings.Contains(got, "issues") {
		t.Errorf("unexpected 'issues' in error string: %s", got)
	}
}

func TestParseError_ValidJSON(t *testing.T) {
	body := []byte(`{"code":"not_found","message":"customer not found","issues":[]}`)
	e := ParseError(404, body)
	if e.StatusCode != 404 {
		t.Errorf("StatusCode: got %d, want 404", e.StatusCode)
	}
	if e.Code != "not_found" {
		t.Errorf("Code: got %q, want not_found", e.Code)
	}
	if e.Message != "customer not found" {
		t.Errorf("Message: got %q, want customer not found", e.Message)
	}
}

func TestParseError_StringIssues(t *testing.T) {
	body := []byte(`{"code":"validation","message":"bad","issues":["field required"]}`)
	e := ParseError(400, body)
	if len(e.Issues) != 1 || e.Issues[0] != "field required" {
		t.Errorf("Issues: got %v, want [field required]", e.Issues)
	}
}

func TestParseError_ObjectIssues(t *testing.T) {
	body := []byte(`{"code":"validation","message":"bad","issues":[{"message":"email invalid"}]}`)
	e := ParseError(400, body)
	if len(e.Issues) != 1 || e.Issues[0] != "email invalid" {
		t.Errorf("Issues: got %v, want [email invalid]", e.Issues)
	}
}

func TestParseError_InvalidJSON(t *testing.T) {
	body := []byte(`not json`)
	e := ParseError(500, body)
	if e.StatusCode != 500 {
		t.Errorf("StatusCode: got %d, want 500", e.StatusCode)
	}
	// Should still return an APIError with empty fields, not panic.
	if e.Code != "" {
		t.Errorf("Code: got %q, want empty", e.Code)
	}
}
