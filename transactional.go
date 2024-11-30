package plunk

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type TransactionalEmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	From    string `json:"from,omitempty"`
	Name    string `json:"name,omitempty"`
}

type ContactInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type EmailRecipient struct {
	Contact ContactInfo `json:"contact"`
	Email   string      `json:"email"`
}

type TransactionalEmailResponse struct {
	Success bool             `json:"success"`
	Emails  []EmailRecipient `json:"emails"`
	Timestamp string         `json:"timestamp"`
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
func (p *Plunk) SendTransactionalEmail(payload TransactionalEmailPayload) (*TransactionalEmailResponse, error) {
	res, err := p.sendTransactionalEmails([]TransactionalEmailPayload{payload})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrEmptyResponse
	}

	return res[0], nil
}

func (p *Plunk) SendMultipleTransactionalEmails(payload []TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	return p.sendTransactionalEmails(payload)
}

func (p *Plunk) sendTransactionalEmails(payload []TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	sem := make(chan bool, 10)
	var wg sync.WaitGroup

	// validate payload
	for _, pl := range payload {
		if pl.To == "" {
			return nil, ErrMissingTo
		}

		if pl.Subject == "" {
			return nil, ErrMissingSubject
		}

		if pl.Body == "" {
			return nil, ErrMissingBody
		}
	}

	result := []*TransactionalEmailResponse{}
	url := p.url(transactionalEmailEndpoint)
	for _, pl := range payload {
		wg.Add(1)
		sem <- true

		go func(pl TransactionalEmailPayload) {
			defer wg.Done()
			defer func() { <-sem }()

			res := &TransactionalEmailResponse{}
			resp, err := p.sendRequest(SendConfig{
				Body:   pl,
				Url:    url,
				Method: http.MethodPost,
			})

			if err != nil {
				return
			}

			defer resp.Body.Close()

			err = decodeResponse(resp, &res)
			if err != nil {
				return
			}

			result = append(result, res)
		}(pl)
	}

	wg.Wait()

	close(sem)

	p.logInfo(fmt.Sprintf("Sent %d transactional emails", len(result)))

	return result, nil
}
