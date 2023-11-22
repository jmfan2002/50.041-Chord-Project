package util

import (
	"encoding/json"
	"io"
)

/*
ReadBody reads the a request or response body into a destination structure
- respBody: the body of a request (request.Body()) or response (response.Body())
- destination: the POINTER to destination (&destination)
*/
func ReadBody(respBody io.ReadCloser, destination interface{}) error {
	defer respBody.Close()

	err := json.NewDecoder(respBody).Decode(&destination)
	if err != nil {
		return err
	}

	return nil
}
