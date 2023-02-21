// Package openapi3 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.1 DO NOT EDIT.
package openapi3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// CreateURL request  with any body
	CreateURLWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateURL(ctx context.Context, body CreateURLJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteURL request
	DeleteURL(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ReadURL request
	ReadURL(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) CreateURLWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateURLRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateURL(ctx context.Context, body CreateURLJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateURLRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteURL(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteURLRequest(c.Server, uRLId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ReadURL(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewReadURLRequest(c.Server, uRLId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewCreateURLRequest calls the generic CreateURL builder with application/json body
func NewCreateURLRequest(server string, body CreateURLJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateURLRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateURLRequestWithBody generates requests for CreateURL with any type of body
func NewCreateURLRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/URLs")
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteURLRequest generates requests for DeleteURL
func NewDeleteURLRequest(server string, uRLId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "URLId", runtime.ParamLocationPath, uRLId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/URLs/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewReadURLRequest generates requests for ReadURL
func NewReadURLRequest(server string, uRLId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "URLId", runtime.ParamLocationPath, uRLId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/URLs/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = operationPath[1:]
	}
	operationURL := url.URL{
		Path: operationPath,
	}

	queryURL := serverURL.ResolveReference(&operationURL)

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// CreateURL request  with any body
	CreateURLWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateURLResponse, error)

	CreateURLWithResponse(ctx context.Context, body CreateURLJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateURLResponse, error)

	// DeleteURL request
	DeleteURLWithResponse(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*DeleteURLResponse, error)

	// ReadURL request
	ReadURLWithResponse(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*ReadURLResponse, error)
}

type CreateURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		URL *URL `json:"URL,omitempty"`
	}
	JSON400 *struct {
		Error *string `json:"error,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r CreateURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON500      *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r DeleteURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ReadURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		URL *URL `json:"URL,omitempty"`
	}
	JSON500 *struct {
		Error *string `json:"error,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ReadURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ReadURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateURLWithBodyWithResponse request with arbitrary body returning *CreateURLResponse
func (c *ClientWithResponses) CreateURLWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateURLResponse, error) {
	rsp, err := c.CreateURLWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateURLResponse(rsp)
}

func (c *ClientWithResponses) CreateURLWithResponse(ctx context.Context, body CreateURLJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateURLResponse, error) {
	rsp, err := c.CreateURL(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateURLResponse(rsp)
}

// DeleteURLWithResponse request returning *DeleteURLResponse
func (c *ClientWithResponses) DeleteURLWithResponse(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*DeleteURLResponse, error) {
	rsp, err := c.DeleteURL(ctx, uRLId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteURLResponse(rsp)
}

// ReadURLWithResponse request returning *ReadURLResponse
func (c *ClientWithResponses) ReadURLWithResponse(ctx context.Context, uRLId string, reqEditors ...RequestEditorFn) (*ReadURLResponse, error) {
	rsp, err := c.ReadURL(ctx, uRLId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseReadURLResponse(rsp)
}

// ParseCreateURLResponse parses an HTTP response from a CreateURLWithResponse call
func ParseCreateURLResponse(rsp *http.Response) (*CreateURLResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &CreateURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			URL *URL `json:"URL,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteURLResponse parses an HTTP response from a DeleteURLWithResponse call
func ParseDeleteURLResponse(rsp *http.Response) (*DeleteURLResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &DeleteURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseReadURLResponse parses an HTTP response from a ReadURLWithResponse call
func ParseReadURLResponse(rsp *http.Response) (*ReadURLResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ReadURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			URL *URL `json:"URL,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest struct {
			Error *string `json:"error,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}
