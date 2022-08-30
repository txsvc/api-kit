package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

	// ErrApiInvocationError indicates an error in an API call
	ErrApiInvocationError = errors.New("api invocation error")
)

// Client - API client encapsulating the http client
type (
	Client struct {
		httpClient *http.Client
		cfg        *settings.DialSettings
		logger     logger.Logger
		userAgent  string
		trace      string
	}
)

func NewClient(ds *settings.DialSettings, logger logger.Logger) (*Client, error) {
	var _ds *settings.DialSettings

	httpClient := NewTransport(logger, http.DefaultTransport)

	// create or clone the settings
	if ds != nil {
		c := ds.Clone()
		_ds = &c
	} else {
		_ds = config.GetConfig().Settings()
		if _ds.Credentials == nil {
			_ds.Credentials = &settings.Credentials{} // just provide something to prevent NPEs further down
		}
	}

	return &Client{
		httpClient: httpClient,
		cfg:        _ds,
		logger:     logger,
		userAgent:  config.GetConfig().Info().UserAgentString(),
		trace:      stdlib.GetString(config.ForceTraceEnv, ""),
	}, nil // FIXME: nothing creates an error here, remove later?
}

// GET is used to request data from the API. No payload, only queries!
func (c *Client) GET(uri string, response interface{}) (int, error) {
	return c.request("GET", fmt.Sprintf("%s%s", c.cfg.Endpoint, uri), nil, response)
}

func (c *Client) POST(uri string, request, response interface{}) (int, error) {
	return c.request("POST", fmt.Sprintf("%s%s", c.cfg.Endpoint, uri), request, response)
}

func (c *Client) PUT(uri string, request, response interface{}) (int, error) {
	return c.request("PUT", fmt.Sprintf("%s%s", c.cfg.Endpoint, uri), request, response)
}

func (c *Client) DELETE(uri string, request, response interface{}) (int, error) {
	return c.request("DELETE", fmt.Sprintf("%s%s", c.cfg.Endpoint, uri), request, response)
}

func (c *Client) request(method, url string, request, response interface{}) (int, error) {
	var req *http.Request

	if request != nil {
		p, err := json.Marshal(&request)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(p))
		if err != nil {
			return http.StatusBadRequest, err
		}
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return http.StatusBadRequest, err
		}
	}

	return c.roundTrip(req, response)
}

func (c *Client) roundTrip(req *http.Request, response interface{}) (int, error) {

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
		return resp.StatusCode, ErrApiInvocationError
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
