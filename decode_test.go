package plunk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestDecodeResponse(t *testing.T) {
	testCases := []struct {
		responseBody string
		expected     TestStruct
	}{
		{
			responseBody: `{"name": "Test", "value": 42}`,
			expected:     TestStruct{Name: "Test", Value: 42},
		},
	}

	for _, tc := range testCases {
		// Create a test server that returns the mocked response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tc.responseBody))
		}))
		defer ts.Close()

		// Send a request to the test server
		resp, err := http.Get(ts.URL)
		assert.Nil(t, err)
		defer resp.Body.Close()

		// Call decodeResponse with the test server's response
		var result TestStruct
		err = decodeResponse(resp, &result)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, tc.expected, result)
	}
}

func TestDecodeStringToMap(t *testing.T) {
	testCases := []struct {
		input    *string
		expected map[string]interface{}
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input: func() *string {
				s := `{"key": "value"}`
				return &s
			}(),
			expected: map[string]interface{}{"key": "value"},
		},
		{
			input: func() *string {
				s := `{"key1": 1, "key2": 2}`
				return &s
			}(),
			expected: map[string]interface{}{"key1": 1.0, "key2": 2.0},
		},
	}

	for _, tc := range testCases {
		result, err := decodeStringToMap(tc.input)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, result)
	}
}

func TestConvertMapToJSONString(t *testing.T) {
	testCases := []struct {
		input    map[string]interface{}
		expected string
	}{
		{
			input:    nil,
			expected: "",
		},
		{
			input:    map[string]interface{}{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			input:    map[string]interface{}{"key1": 1, "key2": 2},
			expected: `{"key1":1,"key2":2}`,
		},
	}

	for i, tc := range testCases {
		result, err := convertMapToJSONString(tc.input)
		assert.Nil(t, err)
		assert.Equal(t, tc.expected, result, "test case %d", i)
	}
}
