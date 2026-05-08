package openrouter

import "net/http"

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}
