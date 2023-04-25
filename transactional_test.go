package plunk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendTransactionalEmail(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	payload := TransactionalEmailPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	res, err := p.SendTransactionalEmail(payload)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Success, true)
}

func TestSendMultipleTransactionalEmails(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	payload := []TransactionalEmailPayload{
		{
			To:      "test@example.com",
			Subject: "Test Subject",
			Body:    "Test Body",
		},
		{
			To:      "test@example.com",
			Subject: "Test Subject 2",
			Body:    "# Test Body 2",
		},
	}

	res, err := p.SendMultipleTransactionalEmails(payload)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res[0].Success, true)
	assert.Equal(t, res[1].Success, true)

	for _, r := range res {
		fmt.Println("Response: ", r)
		assert.Equal(t, r.Success, true)
	}
}

func TestSendTransactionalEmailWithInvalidPayload(t *testing.T) {
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	testCases := []struct {
		payload TransactionalEmailPayload
		err     error
	}{
		{
			payload: TransactionalEmailPayload{
				To:      "",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			err: ErrMissingTo,
		},
		{
			payload: TransactionalEmailPayload{
				To:      "test@example.com",
				Subject: "",
				Body:    "Test Body",
			},
			err: ErrMissingSubject,
		},
		{
			payload: TransactionalEmailPayload{
				To:      "test@example.com",
				Subject: "Test Subject",
				Body:    "",
			},
			err: ErrMissingBody,
		},
	}

	for _, tc := range testCases {
		_, err := p.sendTransactionalEmails([]TransactionalEmailPayload{tc.payload})
		assert.NotNil(t, err)
		assert.Equal(t, err, tc.err)
	}
}
