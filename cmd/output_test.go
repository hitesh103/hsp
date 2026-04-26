package cmd

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestRequestPanelRendering(t *testing.T) {
	fmt.Println("=" + Pad("=", 78))
	fmt.Println("SAMPLE OUTPUT - ASCII TUI Panel Renderer")
	fmt.Println("=" + Pad("=", 78))

	fmt.Println("\n--- Box Functions ---")
	fmt.Println(DrawBox("Request", 60))
	fmt.Println(DrawSection("Section Title", 60))
	fmt.Println(DrawDoubleBox("Double Box", 60))

	fmt.Println("\n--- Request Preview ---")
	req := &RequestBuilder{
		URL:     "https://api.example.com/users/123",
		Method:  "POST",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
		QueryParams: map[string]string{
			"include": "profile",
		},
		Body: `{
  "name": "John Doe",
  "email": "john@example.com"
}`,
	}
	fmt.Println(RenderRequest(req))

	fmt.Println("\n--- Request Preview (Compact) ---")
	fmt.Println(RenderRequestPreview(req))

	fmt.Println("\n--- Response Display ---")
	respHeaders := http.Header{
		"Content-Type":  []string{"application/json"},
		"X-Request-Id": []string{"abc123"},
	}
	respBody := []byte(`{
  "id": 123,
  "name": "John Doe",
  "email": "john@example.com"
}`)
	fmt.Println(RenderResponse(201, "Created", 143*time.Millisecond, respHeaders, respBody))

	fmt.Println("\n--- Helper Functions ---")
	fmt.Println("Truncate('hello world', 8):", Truncate("hello world", 8))
	fmt.Println("Pad('hi', 10):", Pad("hi", 10))
}