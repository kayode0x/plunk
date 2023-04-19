package plunk

type ContactsResponse struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	Data       map[string]interface{} `json:"data"`
}

type CreateContactPayload struct {
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	Data       map[string]interface{} `json:"data"`
}

// Gets the details of a specific contact.
func (p *Plunk) GetContact(id string) (*ContactsResponse, error) {
	return nil, nil
}

// Get a list of all contacts in your Plunk account.
func (p *Plunk) GetContacts() ([]*ContactsResponse, error) {
	return nil, nil
}

// Gets the total number of contacts in your Plunk account.
// Useful for displaying the number of contacts in a dashboard, landing page or other marketing material.
func (p *Plunk) GetContactsCount() (int, error) {
	return 0, nil
}

// Used to create a new contact in your Plunk project without triggering an event
func (p *Plunk) CreateContact(payload *CreateContactPayload) (*ContactsResponse, error) {
	return nil, nil
}

// Updates a contact's subscription status to subscribed.
func (p *Plunk) SubscribeContact(id string) (*ContactsResponse, error) {
	return subOrUnsubscribeContact(id, true)
}

// Updates a contact's subscription status to unsubscribed.
func (p *Plunk) UnsubscribeContact(id string) (*ContactsResponse, error) {
	return subOrUnsubscribeContact(id, false)
}

func subOrUnsubscribeContact(id string, subscribed bool) (*ContactsResponse, error) {
	return nil, nil
}
