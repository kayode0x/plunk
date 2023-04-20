package plunk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Define a struct to represent the JSON error data
type CustomError struct {
	Code    int    `json:"code"`
	Type    string `json:"error"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Plunk Error (Code: %d, Error: %s, Message: %s)", e.Code, e.Type, e.Message)
}

func parseAPIError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var plunkError CustomError
	err = json.Unmarshal(body, &plunkError)
	if err != nil {
		return err
	}

	return &plunkError
}
