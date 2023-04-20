package plunk

import (
	"errors"
	"fmt"
	"net/http"
)

type EventPayload struct {
	Event      string                 `json:"event"`
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"` // When you trigger an event for a contact, they will automatically be subscribed unless you pass subscribed: false along with the event.
	Data       map[string]interface{} `json:"data"`
}

type EventResponse struct {
	Success bool   `json:"success"`
	Contact string `json:"contact"`
}

var (
	ErrMissingEvent = errors.New("missing event")
	ErrMissingEmail = errors.New("missing email")
)

// Triggers an event and creates it if it doesn't exist.
func (p *Plunk) TriggerEvent(payload *EventPayload) (*EventResponse, error) {
	// validate payload
	if payload.Event == "" {
		return nil, ErrMissingEvent
	}

	if payload.Email == "" {
		return nil, ErrMissingEmail
	}

	result := &EventResponse{}
	url := p.url(eventsEndpoint)
	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodPost,
		Body:   payload,
	})

	if err != nil {
		return nil, err
	}

	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	p.logInfo(fmt.Sprintf("Event triggered: %s ", payload.Event))

	return result, nil
}
