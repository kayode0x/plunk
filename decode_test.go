package plunk

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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

	for i, tc := range testCases {
		// Create a test server that returns the mocked response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tc.responseBody))
		}))
		defer ts.Close()

		// Send a request to the test server
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Test case %d: unexpected error during request: %v", i, err)
			continue
		}
		defer resp.Body.Close()

		// Call decodeResponse with the test server's response
		var result TestStruct
		err = decodeResponse(resp, &result)
		if err != nil {
			t.Errorf("Test case %d: unexpected error during decode: %v", i, err)
		}

		// Compare the decoded result with the expected value
		if result != tc.expected {
			t.Errorf("Test case %d: expected %v, got %v", i, tc.expected, result)
		}
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

	for i, tc := range testCases {
		result, err := decodeStringToMap(tc.input)
		if err != nil {
			t.Errorf("Test case %d: unexpected error: %v", i, err)
		}

		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Test case %d: expected %v, got %v", i, tc.expected, result)
		}
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
		if err != nil {
			t.Errorf("Test case %d: unexpected error: %v", i, err)
		}

		if result != tc.expected {
			t.Errorf("Test case %d: expected %v, got %v", i, tc.expected, result)
		}
	}
}
