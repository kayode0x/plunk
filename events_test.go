package plunk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testEvent      = "test_event"
	eventTestEmail = "test@example.com"
)

func TestTriggerEvent(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	assert.Nil(t, err)

	payload := &EventPayload{
		Event: testEvent,
		Email: eventTestEmail,
	}

	// Test valid payload
	resp, err := p.TriggerEvent(payload)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	_, err = p.DeleteEvent(resp.Event)
	assert.Nil(t, err)

	// Test missing event
	payload = &EventPayload{
		Email: eventTestEmail,
	}
	resp, err = p.TriggerEvent(payload)
	assert.Equal(t, ErrMissingEvent, err)
	assert.Nil(t, resp)

	// Test missing email
	payload = &EventPayload{
		Event: testEvent,
	}
	resp, err = p.TriggerEvent(payload)
	assert.Equal(t, ErrMissingEmail, err)
	assert.Nil(t, resp)

	// Test HTTP request error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p.BaseUrl = server.URL
	payload = &EventPayload{
		Event: testEvent,
		Email: eventTestEmail,
	}
	resp, err = p.TriggerEvent(payload)
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	// Test JSON decoding error
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{invalid json}")
	}))
	defer server.Close()

	p.BaseUrl = server.URL
	payload = &EventPayload{
		Event: testEvent,
		Email: eventTestEmail,
	}
	resp, err = p.TriggerEvent(payload)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func TestDeleteEvent(t *testing.T) {
	p, err := NewClient(secretKey, opts)
	assert.Nil(t, err)

	payload := &EventPayload{
		Event: testEvent,
		Email: eventTestEmail,
	}

	resp, err := p.TriggerEvent(payload)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	id := resp.Event

	event, err := p.DeleteEvent(resp.Event)
	assert.Nil(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, id, event.ID)

	// Test valid ID
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/events", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"id":"`+id+`"}`)
	}))
	defer server.Close()

	// Test missing ID
	event, err = p.DeleteEvent("")
	assert.Equal(t, ErrMissingEventID, err)
	assert.Nil(t, event)

	// Test HTTP request error
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p.BaseUrl = server.URL
	event, err = p.DeleteEvent(id)
	assert.NotNil(t, err)
	assert.Nil(t, event)

	// Test JSON decoding error
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{invalid json}")
	}))
	defer server.Close()

	p.BaseUrl = server.URL
	event, err = p.DeleteEvent(id)
	assert.NotNil(t, err)
	assert.Nil(t, event)
}
