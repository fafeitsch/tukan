package tukan

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Returns an error if either the error parameter is not nil or
// the status code is not a "successful" status code.
// Depending on the condition, the error contains more detailed information.
// Otherwise, nil is returned.
func checkResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		data, _ := ioutil.ReadAll(resp.Body)
		msg := string(data)
		return fmt.Errorf("authentication error, status code: %d with message \"%s\" and content \"%s\"", resp.StatusCode, resp.Status, msg)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	return nil
}
