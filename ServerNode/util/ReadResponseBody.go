package util

import (
	"encoding/json"
	"net/http"
)

func ReadResponseBody(resp *http.Response, destination interface{}) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(&destination)
	if err != nil {
		return err
	}

	return nil
}
