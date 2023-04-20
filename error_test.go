package plunk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.Nil(t, err)
		defer resp.Body.Close()

		err = parseAPIError(resp)

		if tc.expectFail {
			assert.NotNil(t, err)
		} else {
			if tc.expected == nil {
				assert.Nil(t, err)
			} else {
				customErr, ok := err.(*CustomError)
				assert.True(t, ok)

				if customErr.Code != tc.expected.Code ||
					customErr.Type != tc.expected.Type ||
					customErr.Message != tc.expected.Message ||
					customErr.Time != tc.expected.Time {
					assert.Fail(t, "Test case %d failed", i)
				}
			}
		}
	}
}
