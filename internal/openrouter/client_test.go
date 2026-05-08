package openrouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamChatParsesSSE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"hel\"}}]}\n\n"))
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"lo\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":2,\"completion_tokens\":3,\"total_tokens\":5}}\n\n"))
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, APIKey: "key", HTTPClient: server.Client()}
	stream, err := client.StreamChat(context.Background(), ChatRequest{Model: "m", Messages: []ChatMessage{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("StreamChat() error = %v", err)
	}
	var content string
	var usage *Usage
	for event := range stream {
		content += event.Content
		if event.Usage != nil {
			usage = event.Usage
		}
	}
	if content != "hello" {
		t.Fatalf("content = %q", content)
	}
	if usage == nil || usage.TotalTokens != 5 {
		t.Fatalf("usage = %#v", usage)
	}
}

func TestStreamChatMapsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "no credit", http.StatusPaymentRequired)
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, APIKey: "key", HTTPClient: server.Client()}
	_, err := client.StreamChat(context.Background(), ChatRequest{Model: "m"})
	if err == nil || !strings.Contains(err.Error(), "API-002") {
		t.Fatalf("StreamChat() error = %v, want API-002", err)
	}
}
