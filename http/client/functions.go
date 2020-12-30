package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Request is the interface used to pass calls for http requests
type Request interface {
}

// Response is the interface used to pass calls for http responses
type Response interface {
}

// Get performs an HTTP get to a specified endpoint and returns the result.
// It accepts an endpoint and an expected struct to unmarshal to
// If the unmarshaling fails it will return a generic response
func (c *Client) Get(endpoint string, input Request, headers map[string]string, response Response) error {
	reqbody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, c.target+endpoint, bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

// Post performs an HTTP post to a specified endpoint and returns the result
// It accepts an endpoint, a struct to post, and an expected struct to unmarshal to.
// If the unmarshaling fails (which assumes missing required data) it will return
// a generic response
func (c *Client) Post(endpoint string, input Request, headers map[string]string, response Response) error {
	reqbody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.target+endpoint, bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

// Patch performs an HTTP patch to a specified endpoint and returns the result
// It accepts an endpoint, a struct to post, and an expected struct to unmarshal to.
// If the unmarshaling fails (which assumes missing required data) it will return
// a generic response
func (c *Client) Patch(endpoint string, input Request, headers map[string]string, response Response) error {
	reqbody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, c.target+endpoint, bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

// PostForm performs an HTTP post to a specified endpoint and returns the result
// It accepts an endpoint, a struct to post, and an expected struct to unmarshal to.
// If the unmarshaling fails (which assumes missing required data) it will return
// a generic response
func (c *Client) PostForm(endpoint string, input *url.Values, headers map[string]string, response Response) error {
	req, err := http.NewRequest(http.MethodPost, c.target+endpoint, strings.NewReader(input.Encode()))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

// Put performs an HTTP put to a specified endpoint and returns the result
// It accepts an endpoint, a struct to post, and an expected struct to unmarshal to
// If the unmarshaling fails, a generic response
func (c *Client) Put(endpoint string, input Request, headers map[string]string, response Response) error {
	reqbody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, c.target+endpoint, bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

// Delete performs an HTTP delete to a specified endpoint and returns the result
// It accepts an endpoint, a struct to post, and an expected struct to unmarshal to
// If the unmarshaling fails, a generic response
func (c *Client) Delete(endpoint string, input Request, headers map[string]string, response Response) error {
	reqbody, err := json.Marshal(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, c.target+endpoint, bytes.NewBuffer(reqbody))
	if err != nil {
		return err
	}

	setHeaders(headers, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err := errorCheck(resp.StatusCode); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

//
func errorCheck(code int) error {
	if code != http.StatusOK {
		return fmt.Errorf("status code %d", code)
	}
	return nil
}

func setHeaders(headers map[string]string, req *http.Request) {
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}
