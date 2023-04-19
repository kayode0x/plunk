package plunk

import (
	"errors"
	"net/http"
	"os"
)

var (
	ErrNoAPIKey = errors.New("no API key provided")
)

const (
	TransactionalEmailEndpoint = "/v1/send"
	EventsEndpoint             = "/v1/track"
	ContactsEndpoint           = "/v1/contacts"
)

type Config struct {
	ApiKey  string
	Client  *http.Client
	BaseUrl string
	Debug   bool
}

func defaultConfig() *Config {
	return &Config{
		ApiKey:  "",
		Client:  http.DefaultClient,
		BaseUrl: "https://api.useplunk.com",
		Debug:   false,
	}
}

type Plunk struct {
	*Config
}

// NewClient returns a new Plunk client.
func NewClient(apiKey string, opts ...func(*Config)) (*Plunk, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	config := defaultConfig()
	config.ApiKey = apiKey

	for _, opt := range opts {
		opt(config)
	}

	return &Plunk{
		Config: config,
	}, nil
}

// NewClientFromEnv returns a new Plunk client using the PLUNK_API_KEY environment variable.
func NewClientFromEnv() (*Plunk, error) {
	apiKey := os.Getenv("PLUNK_API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return NewClient(apiKey)
}
