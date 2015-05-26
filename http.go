package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrInvalidEndpoint is returned when the endpoint is not a valid HTTP URL.
	ErrInvalidEndpoint = errors.New("invalid endpoint")

	// ErrConnectionRefused is returned when the client cannot connect to the given endpoint.
	ErrConnectionRefused = errors.New("cannot connect to server endpoint")
)

type client struct {
	HTTPClient  *http.Client
	endpoint    string
	endpointURL *url.URL
}

func newClient(endpoint string) (*client, error) {
	u, err := parseEndpoint(endpoint)
	if err != nil {
		return nil, ErrInvalidEndpoint
	}
	return &client{
		HTTPClient:  http.DefaultClient,
		endpoint:    endpoint,
		endpointURL: u,
	}, nil
}

func parseEndpoint(endpoint string) (*url.URL, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, ErrInvalidEndpoint
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalidEndpoint
	}
	return u, nil
}

func (c *client) do(method, path string, data interface{}, forceJSON bool, formData url.Values) ([]byte, int, error) {
	var params io.Reader
	if data != nil || forceJSON {
		buf, err := json.Marshal(data)
		if err != nil {
			return nil, -1, err
		}
		params = bytes.NewBuffer(buf)
	}
	if formData != nil {
		params = strings.NewReader(formData.Encode())
	}
	req, err := http.NewRequest(method, c.getURL(path), params)
	if err != nil {
		return nil, -1, err
	}
	if data != nil && formData == nil {
		req.Header.Set("Content-Type", "application/json")
	} else if formData != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if method == "POST" {
		req.Header.Set("Content-Type", "plain/text")
	}
	var resp *http.Response
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, -1, ErrConnectionRefused
		}
		return nil, -1, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, resp.StatusCode, newError(resp.StatusCode, body)
	}
	return body, resp.StatusCode, nil
}

func (c *client) getURL(path string) string {
	urlStr := strings.TrimRight(c.endpointURL.String(), "/")
	return fmt.Sprintf("%s%s", urlStr, path)
}

type Error struct {
	Status  int
	Message string
}

func newError(status int, body []byte) *Error {
	return &Error{Status: status, Message: string(body)}
}

func (e *Error) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.Status, e.Message)
}
