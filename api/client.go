package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/txsvc/apikit"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/internal/settings"
	"github.com/txsvc/apikit/logger"
)

const (
	// NamespacePrefix namespace for the client and CLI
	NamespacePrefix = "/a/v1"

	// format error messages
	MsgStatus = "%s. status: %d"
)

var (
	// ErrMissingCredentials indicates that a credentials are is missing
	ErrMissingCredentials = errors.New("missing credentials")
)

// Client - API client encapsulating the http client
type (
	Client struct {
		httpClient *http.Client
		cfg        *settings.Settings
		logger     logger.Logger
		userAgent  string
		trace      string
	}
)

func NewClient(logger logger.Logger) (*Client, error) {
	httpClient, err := NewHTTPClient(logger, http.DefaultTransport)
	if err != nil {
		return nil, err
	}
	cfg := config.GetSettings()
	if cfg.Credentials == nil {
		return nil, ErrMissingCredentials
	}

	return &Client{
		httpClient: httpClient,
		cfg:        cfg,
		logger:     logger,
		userAgent:  config.UserAgentString(),
		trace:      stdlib.GetString("APIKIT_FORCE_TRACE", ""),
	}, nil
}

func (c *Client) GET(uri string, response interface{}) (int, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.cfg.Endpoint, uri), nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return c.invoke(req, response)
}

func (c *Client) invoke(req *http.Request, response interface{}) (int, error) {

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", c.userAgent)
	if c.cfg.Credentials.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.Credentials.Token)
	}
	if c.trace != "" {
		req.Header.Set("Apikit-Force-Trace", c.trace)
	}

	// perform the request
	resp, err := c.httpClient.Transport.RoundTrip(req)
	if err != nil {
		if resp == nil {
			return http.StatusInternalServerError, err
		}
		return resp.StatusCode, err
	}

	defer resp.Body.Close()

	// anything other than OK, Created, Accepted, NoContent is treated as an error
	if resp.StatusCode > http.StatusNoContent {
		if response != nil {
			// as we expect a response, there might be a StatusObject
			status := StatusObject{}
			err = json.NewDecoder(resp.Body).Decode(&status)
			if err != nil {
				return resp.StatusCode, fmt.Errorf(MsgStatus, err.Error(), resp.StatusCode)
			}
			return status.Status, fmt.Errorf(status.Message)
		}
		return resp.StatusCode, apikit.ErrApiError
	}

	// unmarshal the response if one is expected
	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return resp.StatusCode, nil
}
