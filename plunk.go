package plunk

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

var (
	ErrNoAPIKey = errors.New("no API key provided")
)

const (
	transactionalEmailEndpoint  = "/send"
	eventsEndpoint              = "/track"
	contactsEndpoint            = "/contacts"
	contactsCountEndpoint       = "/contacts/count"
	contactsSubscribeEndpoint   = "/contacts/subscribe"
	contactsUnsubscribeEndpoint = "/contacts/unsubscribe"
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
		BaseUrl: "https://api.useplunk.com/v1",
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

// Append the endpoint to the base URL.
func (p *Plunk) url(endpoint string) string {
	return p.BaseUrl + endpoint
}

// logger
func (p *Plunk) log(level string, a interface{}) {
	if p.Debug {
		fmt.Printf("[%s] %s\n", level, a)
	}
}

func (p *Plunk) logError(a any) {
	p.log("ERROR", a)
}

func (p *Plunk) logInfo(a any) {
	p.log("INFO", a)
}
