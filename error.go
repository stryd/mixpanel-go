package mixpanel

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Body string
	Resp *http.Response
}

func (err *APIError) Error() string {
	return fmt.Sprintf("error returned from mixpanel: code: %d, api: %s, response: %s", err.Resp.StatusCode, err.Resp.Request.URL.Path, err.Body)
}
