package plunk

import (
	"errors"
	"fmt"
	"net/http"
)

type TransactionalEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type TransactionalEmailResponse struct {
	Success bool   `json:"success"`
	Contact string `json:"contact"`
}

var (
	ErrEmptyResponse  = errors.New("empty response")
	ErrMissingTo      = errors.New("missing recipient")
	ErrMissingSubject = errors.New("missing subject")
	ErrMissingBody    = errors.New("missing body or html")
)

// Used to send transactional emails to a single recipient or multiple recipients at once.
// Transactional emails are programmatically sent emails that are considered to be part of your application's workflow.
// This could be a password reset email, a billing email or other non-marketing emails.
//
// # Using Markdown
//
// It is possible to use Markdown when sending a transactional email. Plunk will automatically apply the same styling as the email templates you make in the editor.
// Any email with a body that starts with # will be treated as Markdown.
func (p *Plunk) SendTransactionalEmail(payload *TransactionalEmailPayload) (*TransactionalEmailResponse, error) {
	res, err := p.sendTransactionalEmails([]*TransactionalEmailPayload{payload})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrEmptyResponse
	}

	return res[0], nil
}

func (p *Plunk) SendMultipleTransactionalEmails(payload []*TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	return p.sendTransactionalEmails(payload)
}

func (p *Plunk) sendTransactionalEmails(payload []*TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	// validate payload
	for _, msg := range payload {
		if msg.To == "" {
			return nil, ErrMissingTo
		}

		if msg.Subject == "" {
			return nil, ErrMissingSubject
		}

		if msg.Body == "" {
			return nil, ErrMissingBody
		}
	}

	result := []*TransactionalEmailResponse{}
	url := p.url(transactionalEmailEndpoint)
	resp, err := p.sendRequest(SendConfig{
		Url:    url,
		Body:   payload,
		Method: http.MethodPost,
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = decodeResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	p.logInfo(fmt.Sprintf("Fetched %d transactional emails", len(result)))

	return result, nil
}
