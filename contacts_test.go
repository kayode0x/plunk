package plunk

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var opts = &Config{
	Debug: true,
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var (
	testEmail = "user@example.com"
	secretKey = getEnvVariable("PLUNK_SECRET_KEY")
)

func TestGetContact(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	payload := CreateContactPayload{
		Email: testEmail,
	}

	contact, err := p.CreateContact(payload)
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	contact, err = p.GetContact(contact.ID)
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = p.DeleteContact(contact.ID)
	assert.Nil(t, err)
}

func TestGetContacts(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	contacts, err := p.GetContacts()
	assert.Nil(t, err)
	assert.NotNil(t, contacts)
}

func TestGetContactsCount(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	count, err := p.GetContactsCount()
	assert.Nil(t, err)
	assert.NotNil(t, count)
}

func TestCreateContact(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
	}

	payload := CreateContactPayload{
		Data:       data,
		Subscribed: true,
		Email:      testEmail,
	}

	contact, err := p.CreateContact(payload)
	assert.Nil(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, contact.Email, testEmail)
	assert.Equal(t, contact.Subscribed, true)

	err = contact.ParseData()
	assert.Nil(t, err)
	assert.Equal(t, contact.Data["first_name"], "John")
	assert.Equal(t, contact.Data["last_name"], "Doe")

	_, err = p.DeleteContact(contact.ID)
	assert.Nil(t, err)
}

func TestUpdateContact(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
	}
	payload := CreateContactPayload{
		Data:       data,
		Subscribed: true,
		Email:      testEmail,
	}

	contact, err := p.CreateContact(payload)
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	err = contact.ParseData()
	assert.Nil(t, err)

	newEmail := "user2@example.com"
	newData := map[string]interface{}{
		"first_name": "Jane",
		"last_name":  "Domingo",
	}

	newContactData := &Contact{
		ID:         contact.ID,
		Email:      newEmail,
		Data:       newData,
		Subscribed: false,
	}

	newContact, err := p.UpdateContact(newContactData)
	assert.Nil(t, err)
	assert.NotNil(t, newContact)
	assert.Equal(t, newContact.Email, newEmail)
	assert.Equal(t, newContact.Subscribed, false)

	err = newContact.ParseData()
	assert.Nil(t, err)
	assert.Equal(t, newContact.Data["first_name"], "Jane")
	assert.Equal(t, newContact.Data["last_name"], "Domingo")

	_, err = p.DeleteContact(newContact.ID)
	assert.Nil(t, err)
}

func TestDeleteContact(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	payload := CreateContactPayload{
		Email: testEmail,
	}

	contact, err := p.CreateContact(payload)
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	_, err = p.DeleteContact(contact.ID)
	assert.Nil(t, err)

	_, err = p.DeleteContact(contact.ID)
	assert.NotNil(t, err)

	expectedErr := "Plunk Error (Code: 404, Error: Not Found, Message: That contact was not found)"
	assert.Equal(t, err.Error(), expectedErr)
}

func TestSubOrUnsubscribeContact(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	payload := CreateContactPayload{
		Email: testEmail,
	}

	contact, err := p.CreateContact(payload)
	assert.Nil(t, err)
	assert.NotNil(t, contact)

	tests := []struct {
		name     string
		email    string
		sub      bool
		expected bool
	}{
		{
			name:     "unsubscribe",
			sub:      true,
			expected: false,
		},
		{
			name:     "subscribe",
			sub:      false,
			expected: true,
		},
	}

	for _, test := range tests {
		newContact, err := p.subOrUnsubscribeContact(contact.ID, test.sub)
		assert.Nil(t, err)
		assert.NotNil(t, newContact)
		assert.Equal(t, newContact.Email, testEmail)
		assert.Equal(t, newContact.Subscribed, test.expected)
	}

	_, err = p.DeleteContact(contact.ID)
	assert.Nil(t, err)
}
