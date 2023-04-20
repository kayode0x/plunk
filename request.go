package plunk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Url     string
	Headers map[string]string
}

type SendConfig struct {
	Url    string
	Method string
	Body   interface{}
}

func (p *Plunk) defaultReqConfig() *Request {
	return &Request{
		Url: p.BaseUrl,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", p.ApiKey),
		},
	}
}

func (p *Plunk) sendRequest(config SendConfig) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	url := config.Url
	request := p.defaultReqConfig()
	body, err := json.Marshal(config.Body)
	if err != nil {
		p.logError(fmt.Sprintf("error marshalling body: %s", err.Error()))
		return nil, err
	}

	data := bytes.NewBuffer(body)

	if config.Method == http.MethodGet {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			p.logError(fmt.Sprintf("error creating request: %s", err.Error()))
			return nil, err
		}

		for key, value := range request.Headers {
			req.Header.Add(key, value)
		}

		resp, err = p.Client.Do(req)
		if err != nil {
			p.logError(fmt.Sprintf("error sending request: %s", err.Error()))
			return nil, err
		}
	} else {
		req, err := http.NewRequest(config.Method, url, data)
		if err != nil {
			p.logError(fmt.Sprintf("error creating request: %s", err.Error()))
			return nil, err
		}

		for key, value := range request.Headers {
			req.Header.Add(key, value)
		}

		resp, err = p.Client.Do(req)
		if err != nil {
			p.logError(fmt.Sprintf("error sending request: %s", err.Error()))
			return nil, err
		}
	}

	err = parseAPIError(resp)
	if err != nil {
		return nil, err
	}

	err = checkStatusCode(resp)
	if err != nil {
		p.logError(fmt.Sprintf("error checking status code: %s", err.Error()))
		return nil, err
	}

	p.logInfo(fmt.Sprintf("made %s request to %s, status code: %d", config.Method, url, resp.StatusCode))

	return resp, nil
}

func checkStatusCode(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	return fmt.Errorf("invalid response (%d %s)", resp.StatusCode, resp.Status)
}
