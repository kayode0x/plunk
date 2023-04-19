package plunk

import "errors"

type TransactionalEmailPayload struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Body    string                 `json:"body"`
	Data    map[string]interface{} `json:"data"`
	Html    string                 `json:"html"` // or you can pass in a HTML file
}

type TransactionalEmailResponse struct {
	Success bool   `json:"success"`
	Contact string `json:"contact"`
}

var (
	ErrEmptyResponse = errors.New("empty response")
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
	res, err := sendTransactionalEmails([]*TransactionalEmailPayload{payload})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrEmptyResponse
	}

	return res[0], nil
}

func (p *Plunk) SendMultipleTransactionalEmails(payload []*TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	return sendTransactionalEmails(payload)
}

func sendTransactionalEmails([]*TransactionalEmailPayload) ([]*TransactionalEmailResponse, error) {
	return nil, nil
}
