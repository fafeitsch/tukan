package tukan

import (
	"fmt"
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
		return fmt.Errorf("authentication error, status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status code: %d with message \"%s\"", resp.StatusCode, resp.Status)
	}
	return nil
}
