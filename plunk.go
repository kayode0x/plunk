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
	deleteEventEndpoint         = "/events"
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

func (p *Plunk) defaultConfig() *Config {
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

// New returns a new Plunk client.
func New(apiKey string, c *Config) (*Plunk, error) {
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	p := &Plunk{}
	config := p.defaultConfig()

	if c != nil {
		if c.Client != nil {
			config.Client = c.Client
		}

		if c.BaseUrl != "" {
			config.BaseUrl = c.BaseUrl
		}

		if c.Debug {
			config.Debug = c.Debug
		}
	}

	config.ApiKey = apiKey
	p.Config = config

	return p, nil
}

// NewFromEnv returns a new Plunk client using the PLUNK_API_KEY environment variable.
func NewFromEnv() (*Plunk, error) {
	apiKey := os.Getenv("PLUNK_API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return New(apiKey, nil)
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
