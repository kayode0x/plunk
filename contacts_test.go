package plunk

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var opts = func(c *Config) {
	c.Debug = true
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var testEmail = "user@example.com"
var secretKey = getEnvVariable("PLUNK_SECRET_KEY")

func TestGetContact(t *testing.T) {
	id := getEnvVariable("PLUNK_CONTACT_ID")

	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	contact, err := p.GetContact(id)
	if err != nil {
		t.Errorf("failed to get contact: %s", err.Error())
		t.FailNow()
	}

	if contact == nil {
		t.Errorf("contact is nil")
		t.FailNow()
	}

	if contact.ID != id {
		t.Errorf("contact ID is not correct")
	}
}

func TestGetContacts(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	contacts, err := p.GetContacts()
	if err != nil {
		t.Errorf("failed to get contacts: %s", err.Error())
		t.FailNow()
	}

	if contacts == nil {
		t.Errorf("contacts is nil")
	}
}

func TestGetContactsCount(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	count, err := p.GetContactsCount()
	if err != nil {
		t.Errorf("failed to get contacts count: %s", err.Error())
		t.FailNow()
	}

	// count must be at least 0.
	if count < 0 {
		t.Errorf("contacts count is not correct")
	}

	t.Logf("contacts count: %d", count)
}

func TestCreateContact(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
	}

	payload := &CreateContactPayload{
		Data:       data,
		Subscribed: true,
		Email:      testEmail,
	}

	contact, err := p.CreateContact(payload)
	if err != nil {
		t.Errorf("failed to create contact: %s", err.Error())
		t.FailNow()
	}

	if contact == nil {
		t.Errorf("contact is nil")
		t.FailNow()
	}

	if contact.Email != testEmail {
		t.Errorf("contact email is not correct")
	}

	if contact.Subscribed != true {
		t.Errorf("contact subscribed is not correct")
	}

	err = contact.ParseData()
	if err != nil {
		t.Errorf("failed to parse contact data: %s", err.Error())
		t.FailNow()
	}

	if contact.Data["first_name"] != "John" {
		t.Errorf("contact data is not correct")
	}

	if contact.Data["last_name"] != "Doe" {
		t.Errorf("contact data is not correct")
	}

	_, err = p.DeleteContact(contact.ID)
	if err != nil {
		t.Errorf("failed to delete contact: %s", err.Error())
		t.FailNow()
	}
}

func TestUpdateContact(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	data := map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
	}
	payload := &CreateContactPayload{
		Data:       data,
		Subscribed: true,
		Email:      testEmail,
	}

	contact, err := p.CreateContact(payload)
	if err != nil {
		t.Errorf("failed to create contact: %s", err.Error())
		t.FailNow()
	}

	if contact == nil {
		t.Errorf("contact is nil")
		t.FailNow()
	}

	err = contact.ParseData()
	if err != nil {
		t.Errorf("failed to parse contact data: %s", err.Error())
		t.FailNow()
	}

	fmt.Println("OLD DATA: ", contact.Data)

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
	if err != nil {
		t.Errorf("failed to update contact: %s", err.Error())
		t.FailNow()
	}

	if newContact == nil {
		t.Errorf("contact is nil")
		t.FailNow()
	}

	if newContact.Email != newEmail {
		t.Errorf("contact email is not correct")
	}

	if newContact.Subscribed != false {
		t.Errorf("contact subscribed is not correct")
	}

	err = newContact.ParseData()
	if err != nil {
		t.Errorf("failed to parse contact data: %s", err.Error())
		t.FailNow()
	}

	fmt.Println("NEW DATA: ", newContact.Data)

	if newContact.Data["first_name"] != "Jane" {
		t.Errorf("contact data is not correct")
	}

	if newContact.Data["last_name"] != "Domingo" {
		t.Errorf("contact data is not correct")
	}

	_, err = p.DeleteContact(newContact.ID)
	if err != nil {
		t.Errorf("failed to delete contact: %s", err.Error())
		t.FailNow()
	}
}

func TestDeleteContact(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	payload := &CreateContactPayload{
		Email: testEmail,
	}

	contact, err := p.CreateContact(payload)
	if err != nil {
		t.Errorf("failed to create contact: %s", err.Error())
		t.FailNow()
	}

	if contact == nil {
		t.Errorf("contact is nil")
		t.FailNow()
	}

	_, err = p.DeleteContact(contact.ID)
	if err != nil {
		t.Errorf("failed to delete contact: %s", err.Error())
		t.FailNow()
	}

	_, err = p.DeleteContact(contact.ID)
	if err == nil {
		t.Errorf("delete contact should fail")
		t.FailNow()
	}

	expectedErr := "Plunk Error (Code: 404, Error: Not Found, Message: That contact was not found)"
	if err.Error() != expectedErr {
		t.Errorf("failed to delete contact: %s", err.Error())
		t.FailNow()
	}
}

func TestSubOrUnsubscribeContact(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}

	tests := []struct {
		name     string
		email    string
		sub      bool
		expected bool
	}{
		{
			name:     "subscribe",
			email:    testEmail,
			sub:      true,
			expected: true,
		},
		{
			name:     "unsubscribe",
			email:    testEmail,
			sub:      false,
			expected: false,
		},
	}

	// run the tests in parallel
	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := &CreateContactPayload{
				Email:      test.email,
				Subscribed: test.sub,
			}

			contact, err := p.CreateContact(payload)
			if err != nil {
				t.Errorf("failed to create contact: %s", err.Error())
				t.FailNow()
			}

			if contact == nil {
				t.Errorf("contact is nil")
				t.FailNow()
			}

			contact, err = p.subOrUnsubscribeContact(contact.ID, test.sub)
			if err != nil {
				t.Errorf("failed to subscribe or unsubscribe contact: %s", err.Error())
				t.FailNow()
			}

			if contact == nil {
				t.Errorf("contact is nil")
				t.FailNow()
			}

			if contact.Subscribed != test.expected {
				t.Errorf("contact subscribed is not correct")
			}

			_, err = p.DeleteContact(contact.ID)
			if err != nil {
				t.Errorf("failed to delete contact: %s", err.Error())
				t.FailNow()
			}
		})
	}
}
