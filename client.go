package client

import (
    "net/http"
    "net/url"
    "fmt"
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

func (*Client) do(method string, path string)