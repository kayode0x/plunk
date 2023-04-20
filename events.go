package plunk

import (
	"errors"
	"fmt"
	"net/http"
)

type EventPayload struct {
	Event      string `json:"event"`
	Email      string `json:"email"`
	Subscribed bool   `json:"subscribed"` // When you trigger an event for a contact, they will automatically be subscribed unless you pass subscribed: false along with the event.
}

type EventResponse struct {
	Success bool   `json:"success"`
	Contact string `json:"contact"` // the ID of the contact that was triggered
	Event   string `json:"event"`   // the ID of the event that was triggered
}

type Event struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	ProjectID  string `json:"projectId"`
	CampaignID string `json:"campaignId"`
	TemplateID string `json:"templateId"`
}

var (
	ErrMissingEvent        = errors.New("missing event")
	ErrMissingEmail        = errors.New("missing email")
	ErrMissingEventID      = errors.New("missing event id")
	ErrCouldNotDeleteEvent = errors.New("could not delete event")
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

// Deletes an event.
func (p *Plunk) DeleteEvent(id string) (*Event, error) {
	if id == "" {
		return nil, ErrMissingEventID
	}

	url := p.url(deleteEventEndpoint)
	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodDelete,
		Body:   map[string]string{"id": id},
	})

	if err != nil {
		return nil, err
	}

	result := &Event{}
	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	p.logInfo(fmt.Sprintf("Event deleted: %s ", id))

	return result, nil
}
