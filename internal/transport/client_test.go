package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("default timeout", func(t *testing.T) {
		client := NewClient(0)
		if client.Timeout != DefaultTimeout {
			t.Errorf("expected timeout %v, got %v", DefaultTimeout, client.Timeout)
		}
	})

	t.Run("custom timeout", func(t *testing.T) {
		customTimeout := 30 * time.Second
		client := NewClient(customTimeout)
		if client.Timeout != customTimeout {
			t.Errorf("expected timeout %v, got %v", customTimeout, client.Timeout)
		}
	})
}

func TestCreateJSONRequest(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ctx, cancel := CreateRequestContext(DefaultTimeout)
		defer cancel()

		req, err := CreateJSONRequest(ctx, "POST", "http://example.com", map[string]any{
			"key": "value",
		}, map[string]string{
			"Authorization": "Bearer token",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if req.Method != "POST" {
			t.Errorf("expected method POST, got %s", req.Method)
		}

		if req.URL.String() != "http://example.com" {
			t.Errorf("expected URL http://example.com, got %s", req.URL.String())
		}

		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", req.Header.Get("Content-Type"))
		}

		if req.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept application/json, got %s", req.Header.Get("Accept"))
		}

		if req.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("expected Authorization Bearer token, got %s", req.Header.Get("Authorization"))
		}
	})

	t.Run("nil body", func(t *testing.T) {
		ctx, cancel := CreateRequestContext(DefaultTimeout)
		defer cancel()

		req, err := CreateJSONRequest(ctx, "GET", "http://example.com", nil, nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if req.Header.Get("Content-Type") != "" {
			t.Errorf("expected no Content-Type header, got %s", req.Header.Get("Content-Type"))
		}
	})
}

func TestExecuteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := NewClient(DefaultTimeout)
	ctx, cancel := CreateRequestContext(DefaultTimeout)
	defer cancel()

	req, err := CreateJSONRequest(ctx, "GET", server.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	resp, err := ExecuteRequest(client, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestReadResponseBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	client := NewClient(DefaultTimeout)
	ctx, cancel := CreateRequestContext(DefaultTimeout)
	defer cancel()

	req, err := CreateJSONRequest(ctx, "GET", server.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	resp, err := ExecuteRequest(client, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body, err := ReadResponseBody(resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != "test response" {
		t.Errorf("expected 'test response', got '%s'", string(body))
	}
}

func TestDecodeJSONResponse(t *testing.T) {
	jsonData := `{"key": "value", "number": 42}`
	var result map[string]interface{}

	err := DecodeJSONResponse([]byte(jsonData), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("expected key 'value', got %v", result["key"])
	}

	if result["number"] != float64(42) {
		t.Errorf("expected number 42, got %v", result["number"])
	}
}

func TestCheckStatus(t *testing.T) {
	t.Run("success status", func(t *testing.T) {
		err := CheckStatus(http.StatusOK)
		if err != nil {
			t.Errorf("expected no error for status 200, got %v", err)
		}
	})

	t.Run("error status", func(t *testing.T) {
		err := CheckStatus(http.StatusBadRequest)
		if err == nil {
			t.Fatal("expected error for status 400")
		}

		expected := "HTTP request failed with status 400"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestCreateRequestContext(t *testing.T) {
	t.Run("default timeout", func(t *testing.T) {
		ctx, cancel := CreateRequestContext(0)
		defer cancel()

		if ctx == nil {
			t.Fatal("expected context to be non-nil")
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("expected context to have deadline")
		}

		expectedDeadline := time.Now().Add(DefaultTimeout)
		if deadline.Before(expectedDeadline.Add(-time.Second)) || deadline.After(expectedDeadline.Add(time.Second)) {
			t.Errorf("expected deadline around %v, got %v", expectedDeadline, deadline)
		}
	})

	t.Run("custom timeout", func(t *testing.T) {
		customTimeout := 10 * time.Second
		ctx, cancel := CreateRequestContext(customTimeout)
		defer cancel()

		if ctx == nil {
			t.Fatal("expected context to be non-nil")
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("expected context to have deadline")
		}

		expectedDeadline := time.Now().Add(customTimeout)
		if deadline.Before(expectedDeadline.Add(-time.Second)) || deadline.After(expectedDeadline.Add(time.Second)) {
			t.Errorf("expected deadline around %v, got %v", expectedDeadline, deadline)
		}
	})
}
