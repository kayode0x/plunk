package plunk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultReqConfig(t *testing.T) {
	p := &Plunk{
		&Config{
			BaseUrl: "https://api.plunk.com",
			ApiKey:  "test-api-key",
		},
	}

	reqConfig := p.defaultReqConfig()
	assert.Equal(t, p.BaseUrl, reqConfig.Url)
	assert.Equal(t, "application/json", reqConfig.Headers["Content-Type"])

	expectedAuthHeader := fmt.Sprintf("Bearer %s", p.ApiKey)
	assert.Equal(t, expectedAuthHeader, reqConfig.Headers["Authorization"])
}

func TestSendRequest(t *testing.T) {
	// create a new Plunk object with a mocked http.Client
	p, err := New(secretKey, opts)
	assert.Nil(t, err)

	// create a SendConfig object with a GET method and a mocked response body
	config := SendConfig{
		Url:    p.url(contactsCountEndpoint),
		Method: http.MethodGet,
		Body:   nil,
	}

	resp, err := p.sendRequest(config)
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	var body interface{}
	err = decodeResponse(resp, &body)
	assert.Nil(t, err)
	assert.NotNil(t, body)
}

func TestCheckStatusCode(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "200 OK"}
	err := checkStatusCode(resp)
	assert.Nil(t, err)

	resp = &http.Response{StatusCode: 404, Status: "404 Not Found"}
	err = checkStatusCode(resp)
	assert.NotNil(t, err)

	expectedErrMsg := "invalid response (404 404 Not Found)"
	assert.Equal(t, expectedErrMsg, err.Error())
}
