package client

import (
    "net/http"
    "net/url"
    "fmt"
    "json"
)

type doer interface {
  Do(*http.Request) (*http.Response, error)
}

type Client struct {
  baseURL url.URL
  doer doer
}

type Config struct {
  BaseURL string
}

func New(config *Config) *Client {
  if config == nil || config.BaseURL == "" {
    return &Client{
        baseURL: url.URL{
          Scheme: "http",
          Host: "api.kivaws.org"
        },
      doer: &http.Client{},
  }

  baseURL, err := url.Parse(config.BaseURL)
  if err != nil {
    panic(fmt.Sprintf("cannot parse base URL: %q (%v)", config.BaseURL, err))
  }
  return &Client{
    baseURL: *baseURL,
    doer: &http.Client{},
  }
}

// decode json from http response and return as interface understood
// by caller
func (c *Client) do(method string, path string, query url.Values, body io.Reader, v interface{}) error {
  resp, err := c.raw(method, path, query, body)
  if err != nil {
    return fmt.Errorf("error making request: %s", err)
  }
  decode := json.NewDecoder(resp.Body)
  if err = decode.Decode(&v); err !=nil {
    return fmt.Errorf("cannot decode your dumb response: %s", err)
  }
  return nil
}

// make the actual http request and return the http response
func (c *Client) raw(method string, path string, query url.Values) (*http.Response, error) {

}