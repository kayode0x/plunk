package plunk

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type ContactsResponse struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	DataString *string                `json:"data"`
	Data       map[string]interface{} `json:"-"`
}

type CreateContactPayload struct {
	Email      string                 `json:"email"`
	Subscribed bool                   `json:"subscribed"`
	Data       map[string]interface{} `json:"data"`
}

var (
	ErrorMissingContactID = errors.New("missing contact id")
)

func (r *ContactsResponse) ParseData() error {
	data, err := decodeStringToMap(r.DataString)
	if err != nil {
		return err
	}

	r.Data = data

	return nil
}

// Gets the details of a specific contact.
func (p *Plunk) GetContact(id string) (*ContactsResponse, error) {
	if id == "" {
		return nil, ErrorMissingContactID
	}

	result := &ContactsResponse{}
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

	result.ParseData()
	p.logInfo(fmt.Sprintf("Contact retrieved: %s ", id))

	return result, nil
}

// Get a list of all contacts in your Plunk account.
func (p *Plunk) GetContacts() ([]*ContactsResponse, error) {
	result := []*ContactsResponse{}
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

		go func(contact *ContactsResponse) {
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

	p.logInfo(fmt.Sprintf("Retrieved %d contacts", len(result)))

	return result, nil
}

// Gets the total number of contacts in your Plunk account.
// Useful for displaying the number of contacts in a dashboard, landing page or other marketing material.
func (p *Plunk) GetContactsCount() (int, error) {
	result := 0
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

	p.logInfo(fmt.Sprintf("Retrieved %d contacts", result))

	return result, nil
}

// Used to create a new contact in your Plunk project without triggering an event
func (p *Plunk) CreateContact(payload *CreateContactPayload) (*ContactsResponse, error) {
	result := &ContactsResponse{}
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

	p.logInfo(fmt.Sprintf("Contact created: %s ", payload.Email))

	return result, nil
}

// Updates a contact's subscription status to subscribed.
func (p *Plunk) SubscribeContact(id string) (*ContactsResponse, error) {
	return p.subOrUnsubscribeContact(id, true)
}

// Updates a contact's subscription status to unsubscribed.
func (p *Plunk) UnsubscribeContact(id string) (*ContactsResponse, error) {
	return p.subOrUnsubscribeContact(id, false)
}

func (p *Plunk) subOrUnsubscribeContact(id string, subscribed bool) (*ContactsResponse, error) {
	if id == "" {
		return nil, ErrorMissingContactID
	}

	result := &ContactsResponse{}
	endpoint := p.url(contactsUnsubscribeEndpoint)
	if subscribed {
		endpoint = p.url(contactsSubscribeEndpoint)
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

	p.logInfo(fmt.Sprintf("Contact updated: %s ", id))

	return result, nil
}
