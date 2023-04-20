package plunk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func decodeResponse(resp *http.Response, v interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func decodeStringToMap(str *string) (map[string]interface{}, error) {
	if str == nil {
		return nil, nil
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(*str), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func convertMapToJSONString(m map[string]interface{}) (string, error) {
	if m == nil {
		return "", nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
