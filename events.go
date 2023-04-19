package plunk

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

// Triggers an event and creates it if it doesn't exist.
func (p *Plunk) TriggerEvent(payload *EventPayload) (*EventResponse, error) {
	return nil, nil
}
