package plunk

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Contact struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	DataString *string                `json:"data"`
	Data       map[string]interface{} `json:"-"`
}

type ContactsCountResponse struct {
	Count int `json:"count"`
}

type CreateContactPayload struct {
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	Data       map[string]interface{} `json:"data"`
}

var (
	ErrMissingContactID           = errors.New("missing contact id")
	ErrCouldNotCreateContact      = errors.New("could not create contact")
	ErrCouldNotGetContact         = errors.New("could not get contact")
	ErrCouldNotGetContacts        = errors.New("could not get contacts")
	ErrCouldNotGetCount           = errors.New("could not get count")
	ErrCouldNotSubscribeContact   = errors.New("could not subscribe contact")
	ErrCouldNotUnsubscribeContact = errors.New("could not unsubscribe contact")
	ErrCouldNotDeleteContact      = errors.New("could not delete contact")
	ErrCouldNotUpdateContact      = errors.New("could not update contact")
)

func (r *Contact) ParseData() error {
	data, err := decodeStringToMap(r.DataString)
	if err != nil {
		return err
	}

	r.Data = data

	return nil
}

// Gets the details of a specific contact.
func (p *Plunk) GetContact(id string) (*Contact, error) {
	if id == "" {
		return nil, ErrMissingContactID
	}

	result := &Contact{}
	endpoint := fmt.Sprintf("%s/%s", contactsEndpoint, id)
	url := p.url(endpoint)

	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodGet,
	})

	if err != nil {
		return nil, err
	}

	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrCouldNotGetContact
	}

	result.ParseData()
	p.logInfo(fmt.Sprintf("Contact retrieved: %s ", id))

	return result, nil
}

// Get a list of all contacts in your Plunk account.
func (p *Plunk) GetContacts() ([]*Contact, error) {
	result := []*Contact{}
	url := p.url(contactsEndpoint)

	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodGet,
	})

	if err != nil {
		p.logError(fmt.Sprintf("Could not send request: %s", err.Error()))
		return nil, err
	}

	err = decodeResponse(resp, &result)
	if err != nil {
		p.logError(fmt.Sprintf("Could not decode response: %s", err.Error()))
		return nil, err
	}

	sem := make(chan bool, 10)
	var wg sync.WaitGroup

	for _, contact := range result {
		wg.Add(1)
		sem <- true

		go func(contact *Contact) {
			defer wg.Done()
			defer func() { <-sem }()

			err := contact.ParseData()
			if err != nil {
				p.logError(fmt.Sprintf("Could not parse data: %s", err.Error()))
			}
		}(contact)
	}

	wg.Wait()

	close(sem)

	if result == nil {
		return nil, ErrCouldNotGetContacts
	}

	p.logInfo(fmt.Sprintf("Retrieved %d contacts", len(result)))

	return result, nil
}

// Gets the total number of contacts in your Plunk account.
// Useful for displaying the number of contacts in a dashboard, landing page or other marketing material.
func (p *Plunk) GetContactsCount() (int, error) {
	result := &ContactsCountResponse{}
	url := p.url(contactsCountEndpoint)

	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodGet,
	})

	if err != nil {
		return 0, err
	}

	err = decodeResponse(resp, result)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, ErrCouldNotGetCount
	}

	p.logInfo(fmt.Sprintf("Retrieved %d contacts", result))

	return result.Count, nil
}

// Used to create a new contact in your Plunk project without triggering an event
func (p *Plunk) CreateContact(payload *CreateContactPayload) (*Contact, error) {
	result := &Contact{}
	url := p.url(contactsEndpoint)

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

	if result == nil {
		return nil, errors.New("")
	}

	p.logInfo(fmt.Sprintf("Contact created: %s ", payload.Email))

	return result, nil
}

// Update a contact in your Plunk project.
func (p *Plunk) UpdateContact(c *Contact) (*Contact, error) {
	if c.ID == "" {
		return nil, ErrMissingContactID
	}

	result := &Contact{}
	url := p.url(contactsEndpoint)
	dataString, error := convertMapToJSONString(c.Data)
	if error != nil {
		return nil, error
	}

	payload := &Contact{
		ID:         c.ID,
		Email:      c.Email,
		DataString: &dataString,
		Subscribed: c.Subscribed,
	}

	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Body:   payload,
		Method: http.MethodPut,
	})

	if err != nil {
		return nil, err
	}

	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrCouldNotUpdateContact
	}

	p.logInfo(fmt.Sprintf("Contact updated: %s ", c.ID))

	return result, nil
}

// Delete a contact from your Plunk project.
func (p *Plunk) DeleteContact(id string) (*Contact, error) {
	if id == "" {
		return nil, ErrMissingContactID
	}

	url := p.url(contactsEndpoint)
	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodDelete,
		Body:   map[string]string{"id": id},
	})

	if err != nil {
		return nil, err
	}

	result := &Contact{}
	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, ErrCouldNotDeleteContact
	}

	p.logInfo(fmt.Sprintf("Contact deleted: %s ", id))

	return result, nil
}

// Updates a contact's subscription status to subscribed.
func (p *Plunk) SubscribeContact(id string) (*Contact, error) {
	return p.subOrUnsubscribeContact(id, true)
}

// Updates a contact's subscription status to unsubscribed.
func (p *Plunk) UnsubscribeContact(id string) (*Contact, error) {
	return p.subOrUnsubscribeContact(id, false)
}

func (p *Plunk) subOrUnsubscribeContact(id string, subscribe bool) (*Contact, error) {
	if id == "" {
		return nil, ErrMissingContactID
	}

	result := &Contact{}
	endpoint := contactsUnsubscribeEndpoint
	if subscribe {
		endpoint = contactsSubscribeEndpoint
	}

	url := p.url(endpoint)
	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Method: http.MethodPost,
		Body:   map[string]string{"id": id},
	})

	if err != nil {
		return nil, err
	}

	err = decodeResponse(resp, result)
	if err != nil {
		return nil, err
	}

	if result == nil {
		if subscribe {
			return nil, ErrCouldNotSubscribeContact
		}

		return nil, ErrCouldNotUnsubscribeContact
	}

	p.logInfo(fmt.Sprintf("Contact updated: %s ", id))

	return result, nil
}
