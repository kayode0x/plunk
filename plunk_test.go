package plunk

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var secretKey = getEnvVariable("PLUNK_SECRET_KEY")

func TestGetContacts(t *testing.T) {
	opts := func(c *Config) {
		c.Debug = true
	}

	p, err := NewClient(secretKey, opts)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("failed to create client: %s", err))
	}

	contacts, err := p.GetContacts()
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("failed to get contacts: %s", err))
	}

	assert.NotNil(t, contacts)
}

func TestGetContact(t *testing.T) {
	id := "27cfb323-c73c-4d0b-afcb-652544d7ba2d"
	opts := func(c *Config) {
		c.Debug = true
	}

	p, err := NewClient(secretKey, opts)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("failed to create client: %s", err))
	}

	contact, err := p.GetContact(id)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("failed to get contact: %s", err))
	}

	assert.NotNil(t, contact)
	assert.Equal(t, id, contact.ID)
}
