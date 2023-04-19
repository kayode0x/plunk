package plunk

import (
	"context"
	"errors"
	"net/http"
	"os"
)

var (
	ErrNoAPIKey = errors.New("no API key provided")
)

type Plunk struct {
	ApiKey  string
	Client  *http.Client
	BaseUrl string
}

type Message struct {
	Email string
	Event string
	Data  map[string]interface{}
}

// NewPlunk returns a new Plunk client.
func NewPlunk(apiKey string) *Plunk {
	return &Plunk{
		ApiKey:  apiKey,
		Client:  http.DefaultClient,
		BaseUrl: "https://api.useplunk.com/v1",
	}
}

// NewPlunkFromEnv returns a new Plunk client using the PLUNK_API_KEY environment variable.
func NewPlunkFromEnv() (*Plunk, error) {
	apiKey := os.Getenv("PLUNK_API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return NewPlunk(apiKey), nil
}

// Send sends a message to Plunk.
func (p *Plunk) Send(ctx context.Context, m *Message) error {
	return nil
}

// Add variables to the message.
func (m *Message) AddVariable(key string, value interface{}) error {
	if m.Data == nil {
		m.Data = make(map[string]interface{})
	}

	m.Data[key] = value

	return nil
}
