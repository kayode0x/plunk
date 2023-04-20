package plunk

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseAPIError(t *testing.T) {
	testCases := []struct {
		statusCode   int
		responseBody string
		expected     *CustomError
		expectFail   bool
	}{
		{
			statusCode: http.StatusOK,
			responseBody: `
			{
				"data": "success"
			}`,
			expected:   nil,
			expectFail: false,
		},
		{
			statusCode: http.StatusInternalServerError,
			responseBody: `
			{
				"code": 500,
				"error": "Internal Server Error",
				"message": "Contact already exists",
				"time": 1682004825112
			}`,
			expected: &CustomError{
				Code:    500,
				Type:    "Internal Server Error",
				Message: "Contact already exists",
				Time:    1682004825112,
			},
			expectFail: false,
		},
		{
			expectFail: true,
			statusCode: http.StatusInternalServerError,
			responseBody: `
			{
				"code": 500,
				"error": "Internal Server Error",
				"message": "Contact already exists",
				"time": 1682004825112
			`,
			expected: nil,
		},
	}

	for i, tc := range testCases {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(tc.statusCode)
			w.Write([]byte(tc.responseBody))
		}))
		defer ts.Close()

		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Test case %d: unexpected error during request: %v", i, err)
			continue
		}
		defer resp.Body.Close()

		err = parseAPIError(resp)

		if tc.expectFail {
			if err == nil {
				t.Errorf("Test case %d: expected error, got nil", i)
			}
		} else {
			if tc.expected == nil {
				if err != nil {
					t.Errorf("Test case %d: expected no error, got %v", i, err)
				}
			} else {
				customErr, ok := err.(*CustomError)
				if !ok {
					t.Errorf("Test case %d: expected CustomError, got %T", i, err)
					continue
				}
				if customErr.Code != tc.expected.Code ||
					customErr.Type != tc.expected.Type ||
					customErr.Message != tc.expected.Message ||
					customErr.Time != tc.expected.Time {
					t.Errorf("Test case %d: expected %v, got %v", i, tc.expected, customErr)
				}
			}
		}
	}
}
