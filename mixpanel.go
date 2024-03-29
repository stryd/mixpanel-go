package mixpanel

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// API is a collection of simple methods to interact with the mixpanel API.
	API interface {
		EventsAPI
		// TODO implement more APIs
	}
	Config struct {
		ApiUrl         string
		Token          string
		ServiceAccount *ServiceAccount
		ProjectID      string
	}
	ServiceAccount struct {
		Username string
		Secret   string
	}
)

type mixpanelClient struct {
	eventsClient
}

// Mixpanel struct store the mixpanel endpoint and the project token
type internalClient struct {
	config     Config
	httpClient *http.Client
}

// New creates a new Mixpanel Client with the given options.
func New(options ...Option) API {
	return NewWithClient(http.DefaultClient, options...)
}

func NewWithClient(httpClient *http.Client, options ...Option) API {
	config := Config{
		ApiUrl: "https://api.mixpanel.com",
	}
	for i := range options {
		options[i].Apply(&config)
	}
	c := internalClient{
		httpClient: httpClient,
		config:     config,
	}
	return &mixpanelClient{
		eventsClient: eventsClient{
			c: c,
		},
	}
}

func (c *internalClient) send(req *http.Request) error {
	if c.config.ServiceAccount != nil {
		req.SetBasicAuth(c.config.ServiceAccount.Username, c.config.ServiceAccount.Secret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		body, bodyErr := ioutil.ReadAll(resp.Body)
		if bodyErr != nil {
			return &APIError{Body: fmt.Sprintf("error reading body: %v", bodyErr), Resp: resp}
		}
		return &APIError{Body: string(body), Resp: resp}
	}
	return nil
}
